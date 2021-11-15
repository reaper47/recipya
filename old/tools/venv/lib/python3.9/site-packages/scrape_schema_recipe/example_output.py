from .scrape import load

from sys import version_info

if version_info < (3, 9):
    from importlib_resources import files 
else:
    from importlib.resources import files   # type: ignore

from pathlib import Path
from typing import Dict, List, Tuple, Union


_ex_name_filename = {'irish-coffee': 'bevvy-irish-coffee-2019.html',
                     'google': 'google-recipe-example.html',
                     'green-beans': 'flavorful-green-beans.html',
                     'taco-salad': 'mm-skinny-chicken-taco-salad.html',
                     'tart': 'foodista-british-treacle-tart.html',
                     'tea-cake': 'crumb-lemon-tea-cakes-2019.html',
                     'truffles': 'sweetestkitchen-truffles.html'}

example_names = tuple(_ex_name_filename.keys())


def example_output(name: str,
                   python_objects: Union[bool, List, Tuple] = False,
                   nonstandard_attrs: bool = False,
                   migrate_old_schema: bool = True) -> List[Dict]:
    """
    Example data useful for prototyping and debugging.  Calls the load()
    function.

    Note: the variable example_names is a list of the example names.

    Parameters
    ----------

    name : string
        the name of the example

    python_objects : bool, list, tuple  (optional) (defaults to False)
    nonstandard_attrs : bool, optional (defaults to False)
    migrate_old_schema : bool, optional (defaults to True)

         [Note: refer to load() function for documentation about the optional
          boolean variables]
    Returns
    -------
    list
        a list of dictionaries in the style of schema.org/Recipe JSON-LD
        no results - an empty list will be returned
    """
    if name not in example_names:
        raise(ValueError("no example named '{}'".format(name)))

    return load(files(__package__) / 'test_data' / _ex_name_filename[name],
                python_objects=python_objects,
                nonstandard_attrs=nonstandard_attrs,
                migrate_old_schema=migrate_old_schema)
