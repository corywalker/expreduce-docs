# Welcome to MkDocs

For full documentation visit [mkdocs.org](http://mkdocs.org).

## Commands

* `mkdocs new [dir-name]` - Create a new project.
* `mkdocs serve` - Start the live-reloading docs server.
* `mkdocs build` - Build the documentation site.
* `mkdocs help` - Print this help message.

## Project layout

    mkdocs.yml    # The configuration file.
    docs/
        index.md  # The documentation homepage.
        ...       # Other markdown pages, images and other files.

```wl
(* An example highlighting the features of
   this Pygments plugin for Mathematica *)
lissajous::usage = "An example Lissajous curve.\n" <>
                   "Definition: f(t) = (sin(3t + π/2), sin(t))"
lissajous = {Sin[2^^11 # + 0.005`10 * 1*^2 * Pi], Sin[#]} &;

With[{max = 2 Pi, min = 0},
    ParametricPlot[lissajous[t], {t, min, max}] /. x_Line :> {Dashed, x}
]
```

```wl
def fn():
    pass
> D[Sin[x]/x,x]
In:  D[(Sin[x] * x^-1), x]
Out: ((Cos[x] * x^-1) + (Sin[x] * -1 * x^-2))

> Table[a^2,{a,1,10}]
In:  Table[a^2, {a, 1, 10}]
Out: {1, 4, 9, 16, 25, 36, 49, 64, 81, 100}

> Sum[i, {i, 1, n}]
In:  Sum[i, {i, 1, n}]
Out: (2^-1 * n * (1 + n))

> (2^(-1) * n * (1 + n)) /. n->5
In:  (((2^(1 * -1) * n) * (1 + n))) /. ((n) -> (5))
Out: 15

> Total[Table[i,{i,1,5}]]
test = 5
In[1]:= Total[Table[i, {i, 1, 5}]]
Out[1]= 15
```

    :::python
    def fn():
        pass
    # Code goes here ...
