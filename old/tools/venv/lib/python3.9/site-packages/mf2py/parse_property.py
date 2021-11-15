"""functions to parse the properties of elements"""
from __future__ import unicode_literals, print_function

from .dom_helpers import get_attr, get_img_src_alt, get_textContent, try_urljoin
from .datetime_helpers import normalize_datetime, DATETIME_RE, TIME_RE
from . import value_class_pattern

import sys
import re

if sys.version < '3':
    text_type = unicode
    binary_type = str
else:
    text_type = str
    binary_type = bytes

def text(el, base_url=''):
    """Process p-* properties"""

    # handle value-class-pattern
    prop_value = value_class_pattern.text(el)
    if prop_value is not None:
        return prop_value

    prop_value = get_attr(el, "title", check_name=("abbr", "link"))
    if prop_value is None:
        prop_value = get_attr(el, "value", check_name=("data", "input"))
    if prop_value is None:
        prop_value = get_attr(el, "alt", check_name=("img", "area"))
    if prop_value is None:
        prop_value = get_textContent(el, replace_img=True, base_url=base_url)

    return prop_value


def url(el, dict_class, img_with_alt, base_url=''):
    """Process u-* properties"""

    prop_value = get_attr(el, "href", check_name=("a", "area", "link"))
    if prop_value is None:
        prop_value = get_img_src_alt(el, dict_class, img_with_alt, base_url)
        if prop_value is not None:
            return prop_value
    if prop_value is None:
        prop_value = get_attr(el, "src", check_name=("audio", "video", "source", "iframe"))
    if prop_value is None:
        prop_value = get_attr(el, "poster", check_name="video")
    if prop_value is None:
        prop_value = get_attr(el, "data", check_name="object")

    if prop_value is not None:
        return try_urljoin(base_url, prop_value)

    # handle value-class-pattern
    prop_value = value_class_pattern.text(el)
    if prop_value is not None:
        return prop_value

    prop_value = get_attr(el, "title", check_name="abbr")
    if prop_value is None:
        prop_value = get_attr(el, "value", check_name=("data", "input"))
    if prop_value is None:
        prop_value = get_textContent(el)

    return prop_value


def datetime(el, default_date=None):
    """Process dt-* properties

    Args:
      el (bs4.element.Tag): Tag containing the dt-value

    Returns:
      a tuple (string string): a tuple of two strings, (datetime, date)
    """

    # handle value-class-pattern
    prop_value = value_class_pattern.datetime(el, default_date)
    if prop_value is not None:
        return prop_value

    prop_value = get_attr(el, "datetime", check_name=("time", "ins", "del"))
    if prop_value is None:
        prop_value = get_attr(el, "title", check_name="abbr")
    if prop_value is None:
        prop_value = get_attr(el, "value", check_name=("data", "input"))
    if prop_value is None:
        prop_value = get_textContent(el) 

    # if this is just a time, augment with default date
    match = re.match(TIME_RE + '$', prop_value)
    if match and default_date:
        prop_value = '%s %s' % (default_date, prop_value)
        return normalize_datetime(prop_value), default_date

    # otherwise, treat it as a full date
    match = re.match(DATETIME_RE + '$', prop_value)
    return (normalize_datetime(prop_value, match=match),
            match and match.group('date'),)


def embedded(el, base_url=''):
    """Process e-* properties"""
    return {
        'html': el.decode_contents().strip(),    # secret bs4 method to get innerHTML
        'value': get_textContent(el, replace_img=True, base_url=base_url)
    }
