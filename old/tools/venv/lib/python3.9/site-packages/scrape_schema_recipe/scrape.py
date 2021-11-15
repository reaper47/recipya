#
# Copyright 2019-2021 Micah Cochran
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# internal libraries
from dataclasses import dataclass
import datetime
from pathlib import Path
import sys
# for mypy
from typing import Callable, Dict, IO, List, Optional, Tuple, Union

# external libraries
import extruct
import isodate
import requests


_PACKAGE_PATH = Path(__file__).resolve().parent

# read version from VERSION file
__version__ = (_PACKAGE_PATH / 'VERSION').read_text().strip()


# Follow RFC 7231 sec. 5.5.3
USER_AGENT_STR = f'scrape-schema-recipe/{__version__} requests/{requests.__version__}'


@dataclass
class SSRTypeError(TypeError):
    """Custom error that is raised when the input given is not of the correct type."""
    var_name: str
    object_type: type
    expected_types: str
    
    def __str__(self):
        s = f'{self.var_name} is of type "{self.object_type.__name__}", when expecting one of the following type(s): {self.expected_types}'

        return s


def scrape(location: Union[str, IO[str]],
           python_objects: Union[bool, List, Tuple] = False,
           nonstandard_attrs: bool = False, migrate_old_schema: bool = True,
           user_agent_str: Optional[str] = None) -> List[Dict]:
    """
    Parse data in https://schema.org/Recipe format into a list of dictionaries
    representing the recipe data.

    Parameters
    ----------
    location : string or file-like object
        A url, filename, or text_string of HTML, or a file-like object.

    python_objects : bool, list, tuple  (optional)
        when True it translates certain data types into python objects
          dates into datetime.date, datetimes into datetime.datetimes,
          durations as dateime.timedelta.
        when set to a list or tuple only converts types specified to
          python objects:
          when set to either [dateime.date] or [datetime.datetimes] either will
            convert dates.
          when set to [datetime.timedelta] durations will be converted
        when False no conversion is performed
        (defaults to False)

    nonstandard_attrs : bool, optional
        when True it adds nonstandard (for schema.org/Recipe) attributes to the
        resulting dictionaries, that are outside the specification such as:
            '_format' is either 'json-ld' or 'microdata' (how schema.org/Recipe was encoded into HTML)
            '_source_url' is the source url, when 'url' has already been defined as another value
        (defaults to False)

    migrate_old_schema : bool, optional
        when True it migrates the schema from older version to current version
        (defaults to True)

    user_agent_str : string, optional
        overide the user_agent_string with this value.
        (defaults to None)

    Returns
    -------
    list
        a list of dictionaries in the style of schema.org/Recipe JSON-LD
        no results - an empty list will be returned
    """

    data = {}  # type: Dict[str, List[Dict]]

    if not user_agent_str:
        user_agent_str = USER_AGENT_STR

    # make sure that one and only are defined
    url = None
    if isinstance(location, str):
        # Is this a url?
        if location.startswith(("http://", "https://")):
            return scrape_url(location, python_objects=python_objects,
                              nonstandard_attrs=nonstandard_attrs,
                              user_agent_str=user_agent_str)

        # Is this is is a very long string? Perhaps it has HTML content.
        elif len(location) > 255:
            data = extruct.extract(location)

        # Maybe it is a filename?
        else:
            with open(location) as f:
                data = extruct.extract(f.read())
    elif hasattr(location, 'read'):
        # Assume this is some kind of file-like object that can be read.
        data = extruct.extract(location.read())
    else:
        raise SSRTypeError(var_name="location", 
                           object_type=type(location), 
                           expected_types = "string for a url, filename, or text_string of the HTML, or a file-like object")

    scrapings = _convert_to_scrapings(data, nonstandard_attrs, url=url)

    if migrate_old_schema is True:
        scrapings = _migrate_old_schema(scrapings)

    if python_objects is not False:
        scrapings = _pythonize_objects(scrapings, python_objects)

    return scrapings


def load(fp: Union[str, IO[str], Path],
         python_objects: Union[bool, List, Tuple] = False,
         nonstandard_attrs: bool = False,
         migrate_old_schema: bool = True) -> List[Dict]:
    """load a filename or file object to scrape

    Parameters
    ----------
    fp : string or file-like object
        A file name or a file-like object.

    python_objects : bool, list, tuple  (optional)
        when True it translates certain data types into python objects
          dates into datetime.date, datetimes into datetime.datetimes,
          durations as dateime.timedelta.
        when set to a list or tuple only converts types specified to
          python objects:
          when set to either [dateime.date] or [datetime.datetimes] either will
            convert dates.
          when set to [datetime.timedelta] durations will be converted
        when False no conversion is performed
        (defaults to False)

    nonstandard_attrs : bool, optional
        when True it adds nonstandard (for schema.org/Recipe) attributes to the
        resulting dictionaries, that are outside the specification such as:
            '_format' is either 'json-ld' or 'microdata' (how schema.org/Recipe was encoded into HTML)
            '_source_url' is the source url, when 'url' has already been defined as another value
        (defaults to False)

    migrate_old_schema : bool, optional
        when True it migrates the schema from older version to current version
        (defaults to True)

    Returns
    -------
    list
        a list of dictionaries in the style of schema.org/Recipe JSON-LD
        no results - an empty list will be returned

    """

    data = {}  # type: Dict[str, List[Dict]]

    if isinstance(fp, str):
        with open(fp) as f:
            data = extruct.extract(f.read())
    elif isinstance(fp, Path):
        data = extruct.extract(fp.read_text())
    elif hasattr(fp, 'read'):
        # Assume this is some kind of file-like object that can be read.
        data = extruct.extract(fp.read())
    else:
        raise SSRTypeError(var_name="fp", 
                           object_type=type(fp), 
                           expected_types="a filename, pathlib.Path object, or a file-like object")

    scrapings = _convert_to_scrapings(data, nonstandard_attrs)

    if migrate_old_schema is True:
        scrapings = _migrate_old_schema(scrapings)

    if python_objects is not False:
        scrapings = _pythonize_objects(scrapings, python_objects)

    return scrapings


def loads(string: str, python_objects: Union[bool, List, Tuple] = False,
          nonstandard_attrs: bool = False,
          migrate_old_schema: bool = True) -> List[Dict]:
    """scrapes a string

    Parameters
    ----------
    string : string
        A text string of HTML.

    python_objects : bool, list, tuple  (optional)
        when True it translates certain data types into python objects
          dates into datetime.date, datetimes into datetime.datetimes,
          durations as dateime.timedelta.
        when set to a list or tuple only converts types specified to
          python objects:
          when set to either [dateime.date] or [datetime.datetimes] either will
            convert dates.
          when set to [datetime.timedelta] durations will be converted
        when False no conversion is performed
        (defaults to False)

    nonstandard_attrs : bool, optional
        when True it adds nonstandard (for schema.org/Recipe) attributes to the
        resulting dictionaries, that are outside the specification such as:
            '_format' is either 'json-ld' or 'microdata' (how schema.org/Recipe was encoded into HTML)
            '_source_url' is the source url, when 'url' has already been defined as another value
        (defaults to False)

    migrate_old_schema : bool, optional
        when True it migrates the schema from older version to current version
        (defaults to True)

    Returns
    -------
    list
        a list of dictionaries in the style of schema.org/Recipe JSON-LD
        no results - an empty list will be returned

    """

    if not isinstance(string, str):
        raise SSRTypeError(var_name="string", object_type=type(string), expected_types="string")


    data = {}  # type: Dict[str, List[Dict]]
    data = extruct.extract(string)
    scrapings = _convert_to_scrapings(data, nonstandard_attrs)

    if migrate_old_schema is True:
        scrapings = _migrate_old_schema(scrapings)

    if python_objects is not False:
        scrapings = _pythonize_objects(scrapings, python_objects)

    return scrapings


def scrape_url(url: str, python_objects: Union[bool, List, Tuple] = False,
               nonstandard_attrs: bool = False,
               migrate_old_schema: bool = True,
               user_agent_str: str = None) -> List[Dict]:
    """scrape from a URL

    Parameters
    ----------
    url : string
        A url to download data from and scrape.

    python_objects : bool, list, tuple  (optional)
        when True it translates certain data types into python objects
          dates into datetime.date, datetimes into datetime.datetimes,
          durations as dateime.timedelta.
        when set to a list or tuple only converts types specified to
          python objects:
          when set to either [dateime.date] or [datetime.datetimes] either will
            convert dates.
          when set to [datetime.timedelta] durations will be converted
        when False no conversion is performed
        (defaults to False)

    nonstandard_attrs : bool, optional
        when True it adds nonstandard (for schema.org/Recipe) attributes to the
        resulting dictionaries, that are outside the specification such as:
            '_format' is either 'json-ld' or 'microdata' (how schema.org/Recipe was encoded into HTML)
            '_source_url' is the source url, when 'url' has already been defined as another value
        (defaults to False)

    migrate_old_schema : bool, optional
        when True it migrates the schema from older version to current version
        (defaults to True)

    user_agent_str : string, optional
        overide the user_agent_string with this value.
        (defaults to None)

    Returns
    -------
    list
        a list of dictionaries in the style of schema.org/Recipe JSON-LD
        no results - an empty list will be returned


    """

    if not isinstance(url, str):
        raise SSRTypeError(var_name="url", object_type=type(url), expected_types="string")


    data: Dict[str, List[Dict]] = {}
    if not user_agent_str:
        user_agent_str = USER_AGENT_STR

    r = requests.get(url, headers={"User-Agent": user_agent_str}, timeout=5)
    r.raise_for_status()
    data = extruct.extract(r.text, r.url)
    url = r.url

    scrapings = _convert_to_scrapings(data, nonstandard_attrs, url=url)

    if migrate_old_schema is True:
        scrapings = _migrate_old_schema(scrapings)

    if python_objects is not False:
        scrapings = _pythonize_objects(scrapings, python_objects)

    return scrapings


def _convert_json_ld_recipe(rec: Dict,
                            nonstandard_attrs: bool = False,
                            url: str = None) -> Dict:
    """Helper function for _convert_to_scraping
    for a json-ld record adding extra tags"""
    # not sure if a copy is necessary?
    d = rec.copy()
    if nonstandard_attrs is True:
        d['_format'] = 'json-ld'
    # store the url
    if url:
        if d.get('url') and d.get('url') != url and nonstandard_attrs is True:
            d['_source_url'] = url
        else:
            d['url'] = url
    return d


def _convert_to_scrapings(data: Dict[str, List[Dict]],
                          nonstandard_attrs: bool = False,
                          url: str = None) -> List[Dict]:
    """dectects schema.org/Recipe content in the dictionary and extracts it"""
    out = []
    if data['json-ld'] != []:
        for rec in data['json-ld']:
            if rec.get('@type') == 'Recipe':
                d = _convert_json_ld_recipe(rec, nonstandard_attrs, url)
                out.append(d)

            if rec.get('@context') == 'https://schema.org' and '@graph' in rec.keys():
                # walk the graph
                for subrec in rec['@graph']:
                    if subrec['@type'] == 'Recipe':
                        d = _convert_json_ld_recipe(subrec, nonstandard_attrs, url)
                        out.append(d)

    if data['microdata'] != []:
        for rec in data['microdata']:
            if rec['type'] in ('http://schema.org/Recipe',
                               'https://schema.org/Recipe'):
                d = rec['properties'].copy()
                if nonstandard_attrs is True:
                    d['_format'] = 'microdata'
                # add @context and @type for conversion to the JSON-LD
                # style format
                if rec['type'][:6] == 'https:':
                    d['@context'] = 'https://schema.org'
                else:
                    d['@context'] = 'http://schema.org'
                d['@type'] = 'Recipe'

                # store the url
                if url:
                    if d.get('url') and nonstandard_attrs is True:
                        d['_source_url'] = url
                    else:
                        d['url'] = url

                for key in d.keys():
                    if isinstance(d[key], dict) and 'type' in d[key]:
                        type_ = d[key].pop('type')
                        d[key]['@type'] = type_.split('/')[3]

                out.append(d)

    return out


# properties that will be passed into datetime objects
DATETIME_PROPERTIES = frozenset(['dateCreated', 'dateModified',
                                 'datePublished', 'expires'])
DURATION_PROPERTIES = frozenset(['cookTime', 'performTime', 'prepTime',
                                 'totalTime', 'timeRequired'])


def _parse_determine_date_datetime(s: str) -> Union[datetime.datetime,
                                                    datetime.date]:
    """Parse function parses a date, if time is included it parses as a
    datetime.
    """
    if sys.version_info >= (3, 7):
        # Check if the date includes time.
        if 'T' in s:
            return datetime.datetime.fromisoformat(s)
        else:
            return datetime.date.fromisoformat(s)
    else:
        # Check if the date includes time.
        if 'T' in s:
            return isodate.parse_datetime(s)
        else:
            return isodate.parse_date(s)


# Test if lists/tuples have contain matching items
def _have_matching_items(lst1: Union[bool, List, Tuple],
                         lst2: Union[bool, List, Tuple]):
    if isinstance(lst1, bool):
        return lst1

    if isinstance(lst2, bool):
        return lst2

    s = set(lst1).intersection(lst2)
    return len(s) > 0


def _pythonize_objects(scrapings: List[Dict], python_objects: Union[bool,
                       List, Tuple]) -> List[Dict]:

    if python_objects is False:
        # this really should not be happening
        return scrapings

    # this should work, mypy gives error, this isn't bulletproof code
    if python_objects is True or datetime.timedelta in python_objects:  # type: ignore
        # convert ISO 8601 date times into timedelta
        scrapings = _convert_properties_scrape(scrapings, DURATION_PROPERTIES,
                                               isodate.parse_duration)

    if python_objects is True or _have_matching_items((datetime.date, datetime.datetime), python_objects):
        # convert ISO 8601 date times into datetimes.datetime objects
        scrapings = _convert_properties_scrape(scrapings, DATETIME_PROPERTIES,
                                               _parse_determine_date_datetime)

    return scrapings


def _convert_properties_scrape(recipes: List[Dict], properties: frozenset,
                               function: Callable[[str], Union[datetime.datetime, datetime.date]]) -> List[Dict]:
    for i in range(len(recipes)):
        key_set = set(recipes[i].keys())
        for p in key_set.intersection(properties):
            try:
                recipes[i][p] = function(recipes[i][p])
            except (isodate.ISO8601Error, ValueError, TypeError):
                if recipes[i][p] is None:  # TypeError
                    recipes[i].pop(p)
                # otherwise, it's a parse error, just leave the value as is

    return recipes


def _migrate_old_schema(recipes: List[Dict]) -> List[Dict]:
    """Migrate old schema.org/Recipe version to current schema version."""
    for i in range(len(recipes)):
        # rename 'ingredients' to 'recipeIngredient'
        if 'ingredients' in recipes[i]:
            recipes[i]['recipeIngredient'] = recipes[i].pop('ingredients')

    return recipes
