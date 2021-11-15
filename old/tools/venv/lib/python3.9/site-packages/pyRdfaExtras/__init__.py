# -*- coding: utf-8 -*-
"""
Wrapper around RDFLib's Graph object to provide additional or replacement serializers. The issue is that, in RDFLib 2.X, the turtle and the RDF/XML serialization both have some issues (bugs and ugly output). As a result, the package’s own serializers should be registered and used. On the other hand, in RDFLib 3.X this becomes unnecessary, it is better to keep to the library’s own version. This wrapper provides a subclass of RDFLib’s Graph overriding the serialize method to register, if necessary, a different serializer and use that one.

@summary: Shell around RDLib's Graph
@organization: U{World Wide Web Consortium<http://www.w3.org>}
@author: U{Ivan Herman<a href="http://www.w3.org/People/Ivan/">}
@license: This software is available for use under the U{W3C® SOFTWARE NOTICE AND LICENSE<href="http://www.w3.org/Consortium/Legal/2002/copyright-software-20021231">}
@requires: U{Ordered Dictionary (odict)<http://dev.pocoo.org/hg/sandbox/raw-file/tip/odict.py>}, needed only for the JSON-LD serialization if Python 2.6 or lower is used (Python 2.7 has a built in ordered list module). It is included in the distribution
@requires: U{simplejson package by Bob Ippolito<http://undefined.org/python/#simplejson>}, needed only for the JSON-LD serailization if Python 2.5 or lower is used (Python 2.6 has a json implementation included in the distribution). 
"""

"""
$Id: __init__.py,v 1.4 2013-10-16 11:49:32 ivan Exp $ $Date: 2013-10-16 11:49:32 $

"""

import rdflib
if rdflib.__version__ >= "3.0.0" :
	from rdflib	import Graph
else :
	from rdflib.Graph import Graph
from rdflib	import Namespace

_xml_serializer_name	= "my-rdfxml"
_turtle_serializer_name	= "my-turtle"
_json_serializer_name	= "my-json-ld"

# import rdflib_jsonld
# from rdflib_jsonld.serializer import JsonLDSerializer

try:
    from cStringIO import StringIO
except ImportError:
    from StringIO import StringIO

	
#########################################################################################################
class MyGraph(Graph) :
	"""
	Wrapper around RDFLib's Graph object. The issue is that the serializers in RDFLib are buggy:-(
	
	In RDFLib 2.X both the Turtle and the RDF/XML serializations have issues (bugs and ugly output). In RDFLib 3.X
	the Turtle serialization seems to be fine, but the RDF/XML has problems:-(
	
	This wrapper provides a subclass of RDFLib’s Graph overriding the serialize method to register,
	if necessary, a different serializer and use that one.

	@cvar xml_serializer_registered_2: flag to avoid duplicate registration for RDF/XML for rdflib 2.*
	@type xml_serializer_registered_2: boolean
	@cvar xml_serializer_registered_3: flag to avoid duplicate registration for RDF/XML for rdflib 3.*
	@type xml_serializer_registered_3: boolean
	@cvar json_serializer_registered: flag to avoid duplicate registration for JSON-LD for rdflib 3.*
	@type json_serializer_registered: boolean
	@cvar turtle_serializer_registered_2: flag to avoid duplicate registration for Turtle for rdflib 2.*
	@type turtle_serializer_registered_2: boolean
	"""
	xml_serializer_registered_2		= False
	xml_serializer_registered_3		= False
	turtle_serializer_registered_2	= False
	json_serializer_registered      = False
	
	def __init__(self) :
		Graph.__init__(self)

	def _register_XML_serializer_3(self) :
		"""The default XML Serializer of RDFLib 3.X is buggy, mainly when handling lists. An L{own version<serializers.prettyXMLserializer_3>} is
		registered in RDFlib and used in the rest of the package. 
		"""
		if not MyGraph.xml_serializer_registered_3 :
			from rdflib.plugin import register
			from rdflib.serializer import Serializer
			if rdflib.__version__ > "3.1.0" :
				register(_xml_serializer_name, Serializer,
						 "pyRdfaExtras.serializers.prettyXMLserializer_3_2", "PrettyXMLSerializer")
			else :
				register(_xml_serializer_name, Serializer,
						 "pyRdfaExtras.serializers.prettyXMLserializer_3", "PrettyXMLSerializer")
			MyGraph.xml_serializer_registered_3 = True

	def _register_JSON_serializer_3(self) :
		"""JSON LD serializer 
		"""
		if not MyGraph.json_serializer_registered :
			from rdflib.plugin import register
			from rdflib.serializer import Serializer
			try :
				from rdflib_jsonld.serializer import JsonLDSerializer
				register(_json_serializer_name, Serializer,"rdflib_jsonld.serializer", "JsonLDSerializer")
			except :
				register(_json_serializer_name, Serializer,"pyRdfaExtras.serializers.jsonserializer", "JsonSerializer")
			MyGraph.json_serializer_registered = True

	def _register_XML_serializer_2(self) :
		"""The default XML Serializer of RDFLib 2.X is buggy, mainly when handling lists.
		An L{own version<serializers.prettyXMLserializer>} is
		registered in RDFlib and used in the rest of the package. This is not used for RDFLib 3.X.
		"""
		if not MyGraph.xml_serializer_registered_2 :
			from rdflib.plugin import register
			from rdflib.syntax import serializer, serializers
			register(_xml_serializer_name, serializers.Serializer,
					 "pyRdfaExtras.serializers.prettyXMLserializer", "PrettyXMLSerializer")
			MyGraph.xml_serializer_registered_2 = True

	def _register_Turtle_serializer_2(self) :
		"""The default Turtle Serializers of RDFLib 2.X is buggy and not very nice as far as the output is concerned.
		An L{own version<serializers.TurtleSerializer>} is registered in RDFLib and used in the rest of the package.
		This is not used for RDFLib 3.X.
		"""
		if not MyGraph.turtle_serializer_registered_2 :
			from rdflib.plugin import register
			from rdflib.syntax import serializer, serializers
			register(_turtle_serializer_name, serializers.Serializer,
					 "pyRdfaExtras.serializers.turtleserializer", "TurtleSerializer")
			MyGraph.turtle_serialzier_registered_2 = True
			
	def add(self, t) :
		s,p,o = t
		"""Overriding the Graph's add method to filter out triples with possible None values. It may happen
		in case, for example, a host language is not properly set up for the distiller"""
		if s == None or p == None or o == None :
			return
		else :
			Graph.add(self, (s,p,o))
		
	def serialize(self, format = "xml") :
		"""Overriding the Graph's serialize method to adjust the output format"""
		if rdflib.__version__ >= "3.0.0" :
			# this is the easy case
			if format == "xml" or format == "pretty-xml" :
				self._register_XML_serializer_3()
				return Graph.serialize(self, format=_xml_serializer_name)
			elif format == "json-ld" or format == "json" :
				# The new version of the serialziers in RDFLib 3.2.X require this extra round...
				# I do not have the patience of working out why that is so.
				self._register_JSON_serializer_3()
				stream = StringIO()
				Graph.serialize(self, format=_json_serializer_name, destination = stream, auto_compact = True, indent = 4)
				return stream.getvalue()
			elif format == "nt" :
				return Graph.serialize(self, format="nt")
			elif format == "n3" or format == "turtle" :
				retval =""
				return Graph.serialize(self, format="turtle")
		else :
			if format == "xml" or format == "pretty-xml" :
				self._register_XML_serializer_2()
				return Graph.serialize(self, format=_xml_serializer_name)
			elif format == "nt" :
				return Graph.serialize(self, format="nt")
			elif format == "n3" or format == "turtle" :
				self._register_Turtle_serializer_2()
				return Graph.serialize(self, format=_turtle_serializer_name)


