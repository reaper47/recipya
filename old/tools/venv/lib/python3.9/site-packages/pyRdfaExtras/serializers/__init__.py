# -*- coding: utf-8 -*-
"""
Serializer module with a slightly improved versions of the Turtle and Pretty XML serializers.

Both are, actually, almost verbatim copy of the RDFLib Turtle and pretty XML serializer module, respectively.
For more detailed on how serializers work and are registered by RDFLib, please refer to the RDFLib 
descriptions and source code.

The differences, v.a.v. the original, for the Pretty XML serializer:
 - there were bugs in the original in case an rdf:List was generated with Literal content
 - added a pretty print branch for Seq/Alt/Bag
 - removal of the CDATA sections when generating an XML Literal (it is a cautionary support in RDFLib but, because XML literals can be generated under very controlled circumstances only in the case of RDFa, it is unnecessary). 
 - use the character " instead of ' for namespace declarations (to make it uniform with the way attributes are handled
 - increased the initial depth for nesting (to 8)
	
The differences, v.a.v. the original, for the Turtle serializer:
 - there was a bug in the syntax of the @prefix setting
 - the original Turtle serializer insisted on creating prefixes for all URI references, ending up with a large number of fairly unnecessary prefixes that made the output a bit unreadable. This has been changed so that only the 'registered' prefixes are used, ie, those that the original RDFa source contains
 - the original Turtle had a bug and did not generate the shorhands for lists
 - changed the indentation rules for anonymous ('squared') blank nodes	
	
"""