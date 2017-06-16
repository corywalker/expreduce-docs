# Expreduce

Expreduce implements a language with specialized constructs for term rewriting. It is a neat language for a computer algebra system because it is able to express expression manipulation steps in a form very similar to standard math equations. For example, the product rule in calculus can be expressed as:

```wl
D[a_*b_,x_] := D[a,x]*b + a*D[b,x]
```

Now that the kernel understands the product rule, when it later encounters a pattern matching the above LHS, it will recursively apply the product rule until the expression stabilizes.

The term rewriting system and pattern matching engine is fairly advanced. The computer algebra system at this stage is extremely limited, but simple calculus and algebraic manipulation is certainly supported (see examples below). If you are looking for a more mature computer algebra system, please consider using Mathematica (proprietary) or Mathics (open source, Sympy-backed).

## Source code

Expreduce is [on GitHub](https://github.com/corywalker/expreduce).

# Installation
```
go get github.com/corywalker/expreduce/expreduce
```

# Example
To run the example CAS prompt:

```
cd example
go run calc.go
```

```wl
# go run calc.go

Welcome to Expreduce!

In[1]:= D[Sin[x]/x,x]

Out[1]= ((-1 * x^-2 * Sin[x]) + (Cos[x] * x^-1))

In[2]:= Table[a^2,{a,1,10}]

Out[2]= {1, 4, 9, 16, 25, 36, 49, 64, 81, 100}

In[3]:= Sum[i, {i, 1, n}]

Out[3]= (1/2 * (1 + n) * n)

In[4]:= (2^(-1) * n * (1 + n)) /. n->5

Out[4]= 15

In[5]:= Total[Table[i,{i,1,5}]]

Out[5]= 15

In[6]:= bar[1, foo[a, b]] + bar[2, foo[a, b]] /. bar[amatch_Integer, foo[cmatch__]] + bar[bmatch_Integer, foo[cmatch__]] -> bar[amatch + bmatch, foo[cmatch]]

Out[6]= bar[3, foo[a, b]]
```

# Other projects

Expreduce is indeed very similar to Mathics, a similar term rewriting system that uses Sympy as a backend for CAS operations. I created expreduce for a few reasons. The first is that I wanted to learn everything I could about term rewriting systems. The second is that I believe the syntax implemented in here is better suited for building a computer algebra system than using Python to manipulate expressions (as Sympy, and thus Mathics does). Using a language with first-class support for pattern matching and replacement across expression trees is ideal for writing a computer algebra system. This combined with an optimized core can lead to efficient and informed evaluation without much translation work for the programmer when translating equations to code.

# Current limitations

When the engine applies rules for a given symbol, it tries to match the most "specific" rules first. The current definition of specificity is basic now, but can certainly be improved upon. It works in most cases but I can envision cases where it will be wrong. Right now there is no way to override the order of rule application, but it should be simple to add in the future.

The pattern matching system can be very slow, especially when working with `Orderless` expressions with many terms. This is because correctly matching such terms often involves checking many different permutations of a pattern until one finds a match. My theory right now is that the current matching system is behaving naively and that it can be modified to speed things up.

# Future directions

I'm interested in trying to apply Golang's concurrency paradigms to the evaluation sequence. Some low hanging fruit would be to have parallel computation of mapping pure functions onto Lists or other expressions (computing the derivative of a list expressions). Similarly, supporting automatic threading of Listable functions would be nice (computing sin(x) of a large array). The evaluation of an expression often starts with evaluating each of the parameters at the beginning. This could potentially be made concurrent as well. A more complicated but very interesting application would be to break down the pattern matching engine into concurrent components. We would have to be very careful about side effects here, so we might need to overhaul our scoping constructs or somehow restrict access to the EvalState. Another option would be to create a function that predicts if another function has side effects (is this feasible?). A true prediction could allow the system to fall back to non-concurrent evaluation.

Since there are at least two other replacement engines that implement the same syntax that I know of ([another one here](https://github.com/jyh1/mmaclone)), it would could be useful to decide on a standard link protocol such that the replacement engine is independent of the rules that run on top of it. Another layer of abstraction is the frontend. Really the hierarchy goes as follows: core <- rules <- frontend. It would be nice to see all of these functions factored out and interchangeable.

Also of interest is to build up some formal theory on the rule definitions. There should be some pre-existing literature on this, as term-rewriting is a studied field. Some interesting questions to answer are:

1. Given a set of rules for a symbol (or the universe of rules?), can we find duplicates or reduce the set of rules to the most fundamental ones? Answering such a question would improve the efficiency and clarity of the system.
2. Given a set of rules for a symbol (or the universe of rules?), can we prove that the recursive rewrites will terminate? It is fairly easy to write a rule that loops on itself or in cooperation with other rules. Of course, we should remove from consideration some of the imperative symbols such as `While` and `For`. Even with these imperative functions removed, is this question still the halting problem? Can we restrict our considerations enough such that the problem is not the halting problem? Answering this question would improve the stability of the system.

# Development

To run the tests:
```
cd expreduce
go test
```
