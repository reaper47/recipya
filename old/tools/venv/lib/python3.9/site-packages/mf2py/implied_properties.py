from __future__ import unicode_literals, print_function
from . import mf2_classes
from .dom_helpers import get_attr, get_img_src_alt, get_children, get_textContent, try_urljoin
import sys

if sys.version < '3':
    text_type = unicode
    binary_type = str
else:
    text_type = str
    binary_type = bytes


def name(el, base_url=''):
    """Find an implied name property

    Args:
      el (bs4.element.Tag): a DOM element

    Returns:
      string: the implied name value
    """
    def non_empty(val):
        """If alt or title is empty, we don't want to use it as the implied
        name"""
        return val is not None and val != ''

    # if image or area use alt text if not empty
    prop_value = get_attr(el, "alt", check_name=("img", "area"))
    if non_empty(prop_value):
        return text_type(prop_value)

    # if abbreviation use the title if not empty
    prop_value = get_attr(el, "title", check_name="abbr")
    if non_empty(prop_value):
        return text_type(prop_value)

    # find candidate child or grandchild
    poss_child = None
    children = list(get_children(el))
    if len(children) == 1:
        poss_child = children[0]

        # ignore if mf2 root
        if mf2_classes.root(poss_child.get('class',[])):
            poss_child = None

        # if it is not img, area, abbr then find grandchild
        if poss_child and poss_child.name not in ("img", "area", "abbr"):
            grandchildren = list(get_children(poss_child))
            # if only one grandchild
            if len(grandchildren) == 1:
                poss_child = grandchildren[0]
                # if it is not img, area, abbr or is mf2 root then no possible child
                if poss_child.name not in ("img", "area", "abbr") or mf2_classes.root(poss_child.get('class',[])):
                    poss_child = None

    # if a possible child was found
    if poss_child is not None:
        # use alt if possible child is img or area
        prop_value = get_attr(poss_child, "alt", check_name=("img", "area"))
        if non_empty(prop_value):
            return text_type(prop_value)

        # use title if possible child is abbr
        prop_value = get_attr(poss_child, "title", check_name="abbr")
        if non_empty(prop_value):
            return text_type(prop_value)

    # use text if all else fails
    # replace images with alt but not with src in implied name
    # proposal: https://github.com/microformats/microformats2-parsing/issues/35#issuecomment-393615508
    return get_textContent(el, replace_img=True, img_to_src=False, base_url=base_url)


def photo(el, dict_class, img_with_alt, base_url=''):
    """Find an implied photo property

    Args:
      el (bs4.element.Tag): a DOM element
      dict_class: a python class used as a dictionary (set by the Parser object)
      img_with_alt: a flag to enable experimental parsing of alt attribute with img (set by the Parser object)
      base_url (string): the base URL to use, to reconcile relative URLs

    Returns:
      string or dictionary: the implied photo value or implied photo as a dictionary with alt value
    """

    def get_photo_child(children):
        "take a list of children and finds a valid child for photo property"

        # if element has one image child use source if exists and img is
        # not root class
        poss_imgs = [c for c in children if c.name == 'img']
        if len(poss_imgs) == 1:
            poss_img = poss_imgs[0]
            if not mf2_classes.root(poss_img.get('class', [])):
                return poss_img

        # if element has one object child use data if exists and object is
        # not root class
        poss_objs = [c for c in children if c.name == 'object']
        if len(poss_objs) == 1:
            poss_obj = poss_objs[0]
            if not mf2_classes.root(poss_obj.get('class', [])):
                return poss_obj

    # if element is an img use source if exists
    prop_value = get_img_src_alt(el, dict_class, img_with_alt, base_url)
    if prop_value is not None:
        return prop_value

    # if element is an object use data if exists
    prop_value = get_attr(el, "data", check_name="object")
    if prop_value is not None:
        return text_type(prop_value)

    # find candidate child or grandchild
    poss_child = None
    children = list(get_children(el))

    poss_child = get_photo_child(children)

    # if no possible child found then look for grandchild if only one child which is not not mf2 root
    if poss_child is None and len(children) == 1 and not mf2_classes.root(children[0].get('class', [])): 

        grandchildren = list(get_children(children[0]))
        poss_child = get_photo_child(grandchildren)

    # if a possible child was found parse
    if poss_child is not None:
        # img get src
        prop_value = get_img_src_alt(poss_child, dict_class, img_with_alt, base_url)
        if prop_value is not None:
            return prop_value

        # object get data
        prop_value = get_attr(poss_child, "data", check_name="object")
        if prop_value is not None:
            return text_type(prop_value)


def url(el, base_url=''):
    """Find an implied url property

    Args:
      el (bs4.element.Tag): a DOM element
      base_url (string): the base URL to use, to reconcile relative URLs

    Returns:
      string: the implied url value
    """

    def get_url_child(children):
        "take a list of children and finds a valid child for url property"

        # if element has one <a> child use if not root class
        poss_as = [c for c in children if c.name == 'a']
        if len(poss_as) == 1:
            poss_a = poss_as[0]
            if not mf2_classes.root(poss_a.get('class', [])):
                return poss_a

        # if element has one area child use if not root class
        poss_areas = [c for c in children if c.name == 'area']
        if len(poss_areas) == 1:
            poss_area = poss_areas[0]
            if not mf2_classes.root(poss_area.get('class', [])):
                return poss_area


    # if element is a <a> or area use its href if exists
    prop_value = get_attr(el, "href", check_name=("a", "area"))
    if prop_value is not None:  # an empty href is valid
        return text_type(try_urljoin(base_url, prop_value))

    # find candidate child or grandchild
    poss_child = None
    children = list(get_children(el))

    poss_child = get_url_child(children)

    # if no possible child found then look for grandchild if only one child which is not mf2 root
    if poss_child is None and len(children) == 1 and not mf2_classes.root(children[0].get('class', [])): 

        grandchildren = list(get_children(children[0]))
        poss_child = get_url_child(grandchildren)

    # if a possible child was found parse
    if poss_child is not None:
        prop_value = get_attr(poss_child, "href", check_name=("a", "area"))
        if prop_value is not None:  # an empty href is valid
            return text_type(try_urljoin(base_url, prop_value))
