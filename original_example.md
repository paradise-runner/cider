---
___
 ***
  ***
   ***
_____________________________________
 - - -
 **  * ** * ** * **
-     -      -      -
- - - -
apple_notes_id: x-coredata://88D12C64-8E83-4CF5-BA62-D650A0EABD13/ICNote/p392
---

# CommonMark Comprehensive Example

A complete guide to CommonMark markdown features for testing translation.

## Preliminaries

### Backslash Escapes

\!\"\#\$\%\&\'\(\)\*\+\,\-\.\/\:\;\<\=\>\?\@\[\\\]\^\_\`\{\|\}\~

Escaping special characters: \*not emphasized*, \<br/> not a tag, \[not a link](/foo), \`not code`

Escaped backslash: \\*emphasis*

Backslash at end of line (hard line break):
```
foo\
bar
```

### Entity and Numeric Character References

Named entities: &nbsp; &amp; &copy; &AElig; &Dcaron; &frac34; &HilbertSpace;

Decimal numeric references: &#35; &#1234; &#992;

Hexadecimal numeric references: &#X22; &#XD06; &#xcab;

Invalid entities: &nbsp &x; &#; &#x; &MadeUpEntity;

## Blocks and Inlines

### Thematic Breaks

***

### ATX Headings

# H1 Heading
## H2 Heading
### H3 Heading
#### H4 Heading
##### H5 Heading
###### H6 Heading

# Heading with *emphasis* and \*escaped\*

#                  foo with spaces                     

 ### Indented heading
  ## Two spaces
   # Three spaces

## foo ##
  ###   bar    ###

# foo ##################################

##### foo ##

### foo ###     

### foo ### b

# foo#

### foo \###

## foo #\##

# foo \#

Foo bar
# Heading can interrupt paragraph
Bar foo

## 
#
### ###

------

Foo *bar
baz*
====

  Foo *bar
baz*→
====

Foo

  Foo
-
`

<a title="a lot
---


### Fenced Code Blocks

```
simple code block
```

```javascript
let x = 42;
console.log(x);
```

```python
def hello():
    print("world")
```

~~~
tildes work too
~~~

~~~python
def foo():
    pass
~~~

```````````````````````````````` example
let x = 42;
````````````````````````````````

```  
spaces in fence
```

``` foo
info string
```

``` foo bar
baz
```

``` foo\+bar
info with escape
```


### Link Reference Definitions

[foo]: /url "title"

[foo]: /url 'title'

[foo]: /url (title)

[foo]: <http://example.com?find=\*>

[foo]: /url

[foo]: /url
[foo]: /url2

[foo]: /url
[FOO]: /url2

[ΑΓΩ]: /φου


### Blank Lines

aaa

bbb


> - foo
- bar

>     foo
    bar

> ```
> foo
> ```

> foo
    - bar

foo
> bar


### Lists

- foo
- bar

+ foo
+ bar

* foo
* bar

1. foo
2. bar

1. foo
2. bar
3. baz

-
-
-

1. foo
   - bar
2. baz

1. foo

   bar

2. baz

- foo
- bar
  - baz
    - boo

- foo
- bar
- baz

* a
* b
  + c
  + d
    - e
    - f


## Advanced Examples

### Nested Lists with Mixed Types

1. First item
   - Nested bullet
   - Another bullet
2. Second item
   1. Nested ordered
   2. Another ordered
3. Third item

### Code with Complex Formatting

This is a `code span with *emphasis*` inside.

```python
def function(arg1, arg2):
    """
    Docstring
    """
    return arg1 + arg2
```

### Blockquote with Nested Elements

> # Heading in blockquote
>
> Paragraph in blockquote
>
> - List in blockquote
>   - Nested item
>
> > Nested blockquote
>
> `code in blockquote`

### Complex Links and Images

[Link with **bold** and *italic*](/url "Title")

![Image with **alt** text](https://example.com/image.png "Image Title")

[Link with `code` in text](/url)

### Mixed Inline Elements

This paragraph has **bold**, *italic*, ***bold italic***, `code`, [link](/url), and ![image](url "alt").

Escaped characters: \! \@ \# \$ \% \^ \& \* \( \) \+ \- \= \{ \} \[ \] \| \\ \: \; \" \' \< \> \, \. \? \/

### Tables (Extended Markdown)

| Header 1 | Header 2 |
|----------|----------|
| Cell 1   | Cell 2   |
| Cell 3   | Cell 4   |

### Strikethrough (Extended)

~~strikethrough text~~

### Subscript and Superscript (Extended)

H~2~O

E=mc^2^

### Footnotes (Extended)

This is a footnote[^1].

[^1]: The footnote content.


Foo bar

Bar

- Item 1
***

### Lazy Continuation Lines

> Blockquote line 1
continuation without >

### Emphasis Edge Cases

*foo *bar* baz*

***foo** bar*

*foo **bar baz* bim** qux*

*foo **bar** baz*

**foo *bar* baz**

***foo *bar* baz***

**foo **bar** baz**

___foo__ bar__

### Code Span Edge Cases

`` foo ` bar ``

` `` `

`  ``  `

### Link Edge Cases

[foo](/bar\* "ti\*tle")

[foo]

[foo]: /bar\* "ti\*tle"

### Character Encoding

ℋ ⅆ ∲

Æ Ď ¾

### Control Characters

Null replacement: U+0000 becomes U+FFFD

### Precedence Examples

- `one
- two`

---

## Translation Testing Notes

This file is designed to test markdown translation fidelity by including:

1. **Character escaping** - Backslash escapes and special characters
2. **Entity references** - HTML entities and numeric character references
3. **Block structures** - All heading levels, lists, blockquotes, code blocks
4. **Inline formatting** - Emphasis, strong, code, links, images
5. **Edge cases** - Empty elements, mixed nesting, unusual indentation
6. **Precedence rules** - Cases where block or inline structure takes precedence
7. **Whitespace handling** - Tabs, spaces, blank lines
8. **Nested structures** - Lists within lists, blockquotes with lists, etc.
9. **Special syntax** - Setext headings, reference links, HTML blocks
10. **Internationalization** - Unicode characters and combining marks
