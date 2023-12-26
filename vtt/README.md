# VTT Subtitle Format Functions

## BNF

A VTT file is a number of blocks separated by \n\n
Not super formal BNF below:

```
<block> ::= <headerblock> | <styleblock> | <subtitleblock> | <commentblock>
<headerblock> ::= WEBWTT | WEBVTT - <text>
<styleblock> ::= STYLE {stuff here we dont care about}
<commentblock> ::= NOTE\n<text>
<subtitleblock> ::= <subtitleid>\n<subtitlemeta>\n<subtitle>
<subtitleid> ::= ([a-z] | [A-Z] | [0-9])*\n
<subtitlemeta> ::= <subtitleduration> <subtitlecss>
<subtitleduration> ::= <time> --> <time>
<time> ::= hh:mm:ss.ttt | mm:ss.ttt // h, m, s, t are all digits.
<subtitlecss> ::= {also don't care about this}
<subtitle> ::= <tag><text><tag> // Tag & closetag can be the same as we dont care about the contents.
<tag> ::= \<{dont care}\>
<text> ::= Scan anything up until '\n\n' // Also ignore tags inside. F.e <b>asdf</b> is used to style some stuff.
```
