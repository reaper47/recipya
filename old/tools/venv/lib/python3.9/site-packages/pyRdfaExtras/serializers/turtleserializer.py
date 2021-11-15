"""

Serializer for Turtle. This is a slightly modified version of RDFLib's TurtleSerializer module; the
original version had some bugs (in defining prefixes), and the overall output look has also been slightly improved.

"""

import urlparse
from xml.sax.saxutils 			import escape, quoteattr

from rdflib.BNode 				import BNode
from rdflib.Literal 			import Literal
from rdflib.URIRef 				import URIRef
from rdflib.syntax.xml_names 	import split_uri 

from rdflib.syntax.serializers.RecursiveSerializer import RecursiveSerializer
from rdflib.exceptions import Error

from rdflib import RDF, RDFS

XSD = "http://www.w3.org/2001/XMLSchema#"

SUBJECT = 0
VERB = 1
OBJECT = 2

def _Literal_n3(lit,store):
	language = lit.language
	datatype = lit.datatype
	# Modifying the datatype to use a prefixed value rather than the whole URI
	# One has to be a bit careful, because user defined datatypes are also possible
	# The rest is a copy of the original method took over from Literal
	if datatype :
		datatype = store.namespace_manager.normalizeUri(datatype)

	# TODO: We could also chose quotes based on the quotes appearing in the string, i.e. '"' and "'" ...

	# which is nicer?
	# if lit.find("\"")!=-1 or lit.find("'")!=-1 or lit.find("\n")!=-1:
	if lit.find("\n") != -1:
		# Triple quote this string.
		encoded = lit.replace('\\', '\\\\')
		if lit.find('"""')!=-1: 
			# is this ok?
			encoded=encoded.replace('"""','\\"""')
		if encoded.endswith('"'): encoded=encoded[:-1]+"\\\""
		encoded = '"""%s"""' % encoded
	else: 
		encoded = '"%s"' % lit.replace('\n','\\n').replace('\\', '\\\\').replace('"','\\"')
	if language:
		if datatype:    
			return '%s@%s^^%s' % (encoded, language, datatype)
		else:
			return '%s@%s' % (encoded, language)
	else:
		if datatype:
			return '%s^^%s' % (encoded, datatype)
		else:
			return '%s' % encoded
			

class TurtleSerializer(RecursiveSerializer):
	short_name="turtle"
	indentString = "    "
	def __init__(self, store):
		super(TurtleSerializer, self).__init__(store)
		self.reset()
		self.stream = None
		self.store.namespace_manager.bind("xsd","http://www.w3.org/2001/XMLSchema#")

	def reset(self):
		super(TurtleSerializer, self).reset()
		self._shortNames = {}
		self._started = False
	
	def getQName(self, uri):
		if isinstance(uri, URIRef):
			return self.store.namespace_manager.normalizeUri(uri)
		else :
			return None

	def preprocessTriple(self, triple):
		super(TurtleSerializer, self).preprocessTriple(triple)
		p = triple[1]
		if isinstance(p, BNode):
			self._references[p] = self.refCount(p) +1
			
	def label(self, node):
		qname = self.getQName(node)
		if qname is None:
			if isinstance(node,Literal) :
				return _Literal_n3(node,self.store)
			else :
				return node.n3()
		return qname

	def startDocument(self):
		self._started = True
		ns_list= list(self.store.namespaces())
		ns_list.sort()
		if len(ns_list) == 0:
			return
		
		for prefix, uri in ns_list:
			self.write(self.indent() + '@prefix %s: <%s> .\n' % (prefix, uri))
			
		self.write('\n')

	def endDocument(self):
		pass

	def isValidList(self,l): 
		"""Checks if l is a valid RDF list."""			
		return len([ i for i in self.store.items(l) ]) > 0
		
	def doList(self,l):
		while l:
			item = self.store.value(l, RDF.first)
			if item:
				self.path(item, SUBJECT)
				self.write(' ')
				self.subjectDone(l)
			l = self.store.value(l, RDF.rest)
			
	def p_squared(self, node, position):
		if not isinstance(node, BNode) or node in self._serialized or self.refCount(node) > 1 or position == SUBJECT:
			return False

		if self.isValidList(node): 
			# this is a list
			self.write(' ( ')
			self.depth += 2
			self.doList(node)
			self.write(')')
			self.depth -= 2			
		else :
			self.subjectDone(node)
			self.depth += 2
			self.write('\n' + self.indent() + ' [')
			self.predicateList(node)
			self.write('\n' + self.indent() + ' ]')
			self.depth -= 2
		return True

	def p_default(self, node, ignore):
		if ignore == SUBJECT :
			self.write(self.label(node))
		else :
			self.write(" " + self.label(node))
		return True
	
	def path(self, node, position):
		if not (self.p_squared(node, position) or self.p_default(node, position)):
			raise Error("Cannot serialize node '%s'" % (node, ))

	def verb(self, node):
		if node == RDF.type:
			self.write(' a')
		else:
			self.path(node, VERB)
	
	def objectList(self, predicate, objects):
		num = len(objects)
		if num == 0:
			return
		elif num > 3 :
			self.write('\n' + self.indent(2))

		self.path(objects[0], OBJECT)
		for obj in objects[1:]:
			if num > 3 :
				self.write(',\n' + self.indent(2))
			else :
				self.write(",")
			self.path(obj, OBJECT)

	def predicateList(self, subject):
		properties = self.buildPredicateHash(subject)
		propList = self.sortProperties(properties)
		if len(propList) == 0:
			return

		self.verb(propList[0])
		self.objectList(propList[0],properties[propList[0]])
		for predicate in propList[1:]:
			self.write(' ;\n'+self.indent(1))
			self.verb(predicate)
			self.objectList(predicate,properties[predicate])

	def s_squared(self, subject):
		if (self.refCount(subject) > 0) or not isinstance(subject, BNode) :
			return False
		self.write('\n' + self.indent() + " []")
		self.depth += 1
		self.predicateList(subject)
		self.write(' .')
		self.depth -= 1
		#self.write('].')
		return True

	def s_default(self, subject):
		self.write('\n' + self.indent())
		self.path(subject, SUBJECT)
		self.predicateList(subject)
		self.write(' . ')
		return True
	
	def statement(self, subject):
		self.subjectDone(subject)
		if not self.s_squared(subject):
			self.s_default(subject)
			
	def serialize(self, stream, base=None, encoding=None, **args):
		self.reset()
		self.stream = stream
		self.base=base
		
		self.preprocess()
		subjects_list = self.orderSubjects()

		self.startDocument()

		firstTime = True
		for subject in subjects_list:
			if not self.isDone(subject):
				if firstTime:
					firstTime = False
				else:
					self.write('\n')
				self.statement(subject)
		
		self.endDocument()
