# -*- coding: utf-8 -*-
import re

import lxml
import lxml.etree
from lxml.html.clean import Cleaner


NEWLINE_TAGS = frozenset([
    'article', 'aside', 'br', 'dd', 'details', 'div', 'dt', 'fieldset',
    'figcaption', 'footer', 'form', 'header', 'hr', 'legend', 'li', 'main',
    'nav', 'table', 'tr'
])
DOUBLE_NEWLINE_TAGS = frozenset([
    'blockquote', 'dl', 'figure', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'ol',
    'p', 'pre', 'title', 'ul'
])

cleaner = Cleaner(
    scripts=True,
    javascript=False,  # onclick attributes are fine
    comments=True,
    style=True,
    links=True,
    meta=True,
    page_structure=False,  # <title> may be nice to have
    processing_instructions=True,
    embedded=True,
    frames=True,
    forms=False,  # keep forms
    annoying_tags=False,
    remove_unknown_tags=False,
    safe_attrs_only=False,
)


def _cleaned_html_tree(html):
    if isinstance(html, lxml.html.HtmlElement):
        tree = html
    else:
        tree = parse_html(html)

    # we need this as https://bugs.launchpad.net/lxml/+bug/1838497
    try:
        cleaned = cleaner.clean_html(tree)
    except AssertionError:
        cleaned = tree

    return cleaned


def parse_html(html):
    """ Create an lxml.html.HtmlElement from a string with html.
    XXX: mostly copy-pasted from parsel.selector.create_root_node
    """
    body = html.strip().replace('\x00', '').encode('utf8') or b'<html/>'
    parser = lxml.html.HTMLParser(recover=True, encoding='utf8')
    root = lxml.etree.fromstring(body, parser=parser)
    if root is None:
        root = lxml.etree.fromstring(b'<html/>', parser=parser)
    return root


_whitespace = re.compile(r'\s+')
_has_trailing_whitespace = re.compile(r'\s$').search
_has_punct_after = re.compile(r'^[,:;.!?")]').search
_has_open_bracket_before = re.compile(r'\($').search


def _normalize_whitespace(text):
    return _whitespace.sub(' ', text.strip())


def etree_to_text(tree,
                  guess_punct_space=True,
                  guess_layout=True,
                  newline_tags=NEWLINE_TAGS,
                  double_newline_tags=DOUBLE_NEWLINE_TAGS):
    """
    Convert a html tree to text. Tree should be cleaned with
    ``html_text.html_text.cleaner.clean_html`` before passing to this
    function.

    See html_text.extract_text docstring for description of the
    approach and options.
    """
    chunks = []

    _NEWLINE = object()
    _DOUBLE_NEWLINE = object()

    class Context:
        """ workaround for missing `nonlocal` in Python 2 """
        # _NEWLINE, _DOUBLE_NEWLINE or content of the previous chunk (str)
        prev = _DOUBLE_NEWLINE

    def should_add_space(text, prev):
        """ Return True if extra whitespace should be added before text """
        if prev in {_NEWLINE, _DOUBLE_NEWLINE}:
            return False
        if not guess_punct_space:
            return True
        if not _has_trailing_whitespace(prev):
            if _has_punct_after(text) or _has_open_bracket_before(prev):
                return False
        return True

    def get_space_between(text, prev):
        if not text:
            return ' '
        return ' ' if should_add_space(text, prev) else ''

    def add_newlines(tag, context):
        if not guess_layout:
            return
        prev = context.prev
        if prev is _DOUBLE_NEWLINE:  # don't output more than 1 blank line
            return
        if tag in double_newline_tags:
            context.prev = _DOUBLE_NEWLINE
            chunks.append('\n' if prev is _NEWLINE else '\n\n')
        elif tag in newline_tags:
            context.prev = _NEWLINE
            if prev is not _NEWLINE:
                chunks.append('\n')

    def add_text(text_content, context):
        text = _normalize_whitespace(text_content) if text_content else ''
        if not text:
            return
        space = get_space_between(text, context.prev)
        chunks.extend([space, text])
        context.prev = text_content

    def traverse_text_fragments(tree, context, handle_tail=True):
        """ Extract text from the ``tree``: fill ``chunks`` variable """
        add_newlines(tree.tag, context)
        add_text(tree.text, context)
        for child in tree:
            traverse_text_fragments(child, context)
        add_newlines(tree.tag, context)
        if handle_tail:
            add_text(tree.tail, context)

    traverse_text_fragments(tree, context=Context(), handle_tail=False)
    return ''.join(chunks).strip()


def selector_to_text(sel, guess_punct_space=True, guess_layout=True):
    """ Convert a cleaned parsel.Selector to text.
    See html_text.extract_text docstring for description of the approach
    and options.
    """
    import parsel
    if isinstance(sel, parsel.SelectorList):
        # if selecting a specific xpath
        text = []
        for s in sel:
            extracted = etree_to_text(
                s.root,
                guess_punct_space=guess_punct_space,
                guess_layout=guess_layout)
            if extracted:
                text.append(extracted)
        return ' '.join(text)
    else:
        return etree_to_text(
            sel.root,
            guess_punct_space=guess_punct_space,
            guess_layout=guess_layout)


def cleaned_selector(html):
    """ Clean parsel.selector.
    """
    import parsel
    try:
        tree = _cleaned_html_tree(html)
        sel = parsel.Selector(root=tree, type='html')
    except (lxml.etree.XMLSyntaxError,
            lxml.etree.ParseError,
            lxml.etree.ParserError,
            UnicodeEncodeError):
        # likely plain text
        sel = parsel.Selector(html)
    return sel


def extract_text(html,
                 guess_punct_space=True,
                 guess_layout=True,
                 newline_tags=NEWLINE_TAGS,
                 double_newline_tags=DOUBLE_NEWLINE_TAGS):
    """
    Convert html to text, cleaning invisible content such as styles.

    Almost the same as normalize-space xpath, but this also
    adds spaces between inline elements (like <span>) which are
    often used as block elements in html markup, and adds appropriate
    newlines to make output better formatted.

    html should be a unicode string or an already parsed lxml.html element.

    ``html_text.etree_to_text`` is a lower-level function which only accepts
    an already parsed lxml.html Element, and is not doing html cleaning itself.

    When guess_punct_space is True (default), no extra whitespace is added
    for punctuation. This has a slight (around 10%) performance overhead
    and is just a heuristic.

    When guess_layout is True (default), a newline is added
    before and after ``newline_tags`` and two newlines are added before
    and after ``double_newline_tags``. This heuristic makes the extracted
    text more similar to how it is rendered in the browser.

    Default newline and double newline tags can be found in
    `html_text.NEWLINE_TAGS` and `html_text.DOUBLE_NEWLINE_TAGS`.
    """
    if html is None:
        return ''
    cleaned = _cleaned_html_tree(html)
    return etree_to_text(
        cleaned,
        guess_punct_space=guess_punct_space,
        guess_layout=guess_layout,
        newline_tags=newline_tags,
        double_newline_tags=double_newline_tags,
    )
