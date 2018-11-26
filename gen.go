package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"github.com/corywalker/expreduce/expreduce"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func writeMainIndex(fn string) {
	// For more granular writes, open a file for writing.
	os.MkdirAll(path.Dir(fn), os.ModePerm)
	f, err := os.Create(fn)
	check(err)

	// It's idiomatic to defer a `Close` immediately
	// after opening a file.
	defer f.Close()

	f.WriteString("#Expreduce documentation\n")

	f.Sync()
	fmt.Printf("Finished writing %v.\n", fn)
}

func defNamePrint(name string) string {
	//return strings.Replace(name, "`", "\\`", -1)
	return name
}

func defNameFile(name string) string {
	//return strings.ToLower(strings.Replace(name, "`", "_", -1))
	return strings.ToLower(name)
}

func writeCategoryIndex(fn string, defSet expreduce.NamedDefSet) {
	// For more granular writes, open a file for writing.
	os.MkdirAll(path.Dir(fn), os.ModePerm)
	f, err := os.Create(fn)
	check(err)

	// It's idiomatic to defer a `Close` immediately
	// after opening a file.
	defer f.Close()

	f.WriteString(fmt.Sprintf("#%v documentation\n", defSet.Name))

	for _, def := range defSet.Defs {
		if def.OmitDocumentation{
			continue
		}
		f.WriteString(fmt.Sprintf("[%v](%v.md)\n\n", defNamePrint(def.Name), defNameFile(def.Name)))
	}

	f.Sync()
	fmt.Printf("Finished writing %v.\n", fn)
}

func renderUsage(f *os.File, def expreduce.Definition, es *expreduce.EvalState) {
	if len(def.Usage) > 0 {
		f.WriteString(fmt.Sprintf("%v\n\n", def.Usage))
	}
	attrLookup := fmt.Sprintf("Attributes[%s]", defNamePrint(def.Name))
	attrs := expreduce.EasyRun(attrLookup, es)
	f.WriteString(fmt.Sprintf("`%v := %v`\n\n", attrLookup, attrs))
}

func renderExamples(f *os.File, category string, examples []expreduce.TestInstruction, es *expreduce.EvalState) {
	f.WriteString(fmt.Sprintf("##%v\n\n", category))
	count := 1
	for _, ti := range examples {
		comment, isComment := ti.(*expreduce.TestComment)
		if isComment {
			f.WriteString(fmt.Sprintf("%v\n", comment.Comment))
			continue
		}

		inStr, outStr := "", ""
		sameTest, isSameTest := ti.(*expreduce.SameTest)
		if isSameTest {
			inStr = sameTest.In
			outStr = sameTest.Out
		}
		toStringParams := expreduce.ActualStringFormArgsFull("InputForm", es)
		sameTestEx, isSameTestEx := ti.(*expreduce.SameTestEx)
		if isSameTestEx {
			inStr = strings.Replace(sameTestEx.In.StringForm(toStringParams), "Private`", "", -1)
			outStr = strings.Replace(sameTestEx.Out.StringForm(toStringParams), "Private`", "", -1)
		}
		stringTest, isStringTest := ti.(*expreduce.StringTest)
		if isStringTest {
			inStr = stringTest.In
			outStr = stringTest.Out
		}
		exampleOnlyTest, isExampleOnlyInstruction := ti.(*expreduce.ExampleOnlyInstruction)
		if isExampleOnlyInstruction {
			inStr = exampleOnlyTest.In
			outStr = exampleOnlyTest.Out
		}
		if len(inStr) > 0 {
			f.WriteString("```wl\n")
			f.WriteString(fmt.Sprintf("In[%d]:= %v\n", count, inStr))
			f.WriteString(fmt.Sprintf("Out[%d]= %v\n", count, outStr))
			f.WriteString("```\n")
			count += 1
			continue
		}
	}
}

func writeSymbol(fn string, defSet expreduce.NamedDefSet, def expreduce.Definition, es *expreduce.EvalState) {
	// For more granular writes, open a file for writing.
	os.MkdirAll(path.Dir(fn), os.ModePerm)
	f, err := os.Create(fn)
	check(err)

	// It's idiomatic to defer a `Close` immediately
	// after opening a file.
	defer f.Close()

	f.WriteString(fmt.Sprintf("#%v\n\n", defNamePrint(def.Name)))

	renderUsage(f, def, es)

	if len(def.Details) > 0 {
		f.WriteString(fmt.Sprintf("##Details\n\n%v\n\n", def.Details))
	}

	if len(def.SimpleExamples) > 0 {
		renderExamples(f, "Simple examples", def.SimpleExamples, es)
	}

	if len(def.FurtherExamples) > 0 {
		renderExamples(f, "Further examples", def.FurtherExamples, es)
	}

	f.Sync()
	fmt.Printf("Finished writing %v.\n", fn)
}

func main() {
	var docs_location = flag.String("docs_location", "./doc_source", "Location of the docs directory.")
	flag.Parse()

	fmt.Printf("Generating documentation.\n")
	es := expreduce.NewEvalState()

	ymlFn := "mkdocs.yml"
	f, err := os.Create(ymlFn)
	check(err)

	// It's idiomatic to defer a `Close` immediately
	// after opening a file.
	defer f.Close()

	// Generate top level configuration.
	f.WriteString("site_name: Expreduce\n\n")
	f.WriteString("docs_dir: 'doc_source'\n")
	f.WriteString("site_dir: 'docs'\n\n")
	f.WriteString("pages:\n")
	f.WriteString("- Home: 'index.md'\n")
	writeMainIndex(path.Join(*docs_location, "index.md"))
	f.WriteString("- Language reference:\n")

	// Generate module-specific documentation.
	defSets := expreduce.GetAllDefinitions()
	for _, defSet := range defSets {
		categoryFn := fmt.Sprintf("builtin/%s/index.md", defSet.Name)
		writeCategoryIndex(path.Join(*docs_location, categoryFn), defSet)
		categoryDef := fmt.Sprintf(
			"    - '%s': '%s'\n",
			defSet.Name,
			categoryFn,
		)
		f.WriteString(categoryDef)

		for _, def := range defSet.Defs {
			if def.OmitDocumentation{
				continue
			}
			def.AnnotateWithDynamic(es)
			symbolFn := fmt.Sprintf(
				"builtin/%s/%s.md",
				defSet.Name,
				defNameFile(def.Name),
			)
			writeSymbol(path.Join(*docs_location, symbolFn), defSet, def, es)
			symbolDef := fmt.Sprintf(
				"    - '%s ': '%s'\n",
				defNamePrint(def.Name),
				symbolFn,
			)
			f.WriteString(symbolDef)
		}
	}

	// Write remaining configuration.
	f.WriteString("\ntheme: readthedocs\n")
	f.WriteString("theme_dir: 'material'\n")
	f.WriteString("\n")
	f.WriteString("repo_name: 'GitHub'\n")
	f.WriteString("repo_url: 'https://github.com/corywalker/expreduce'\n")
	f.WriteString("\n")
	f.WriteString("extra:\n")
	//f.WriteString("  version: '0.1.0'\n")
	f.WriteString("  logo: 'assets/images/logo.png'\n")
	f.WriteString("  palette:\n")
	f.WriteString("    primary: 'red'\n")
	f.WriteString("    accent: 'light blue'\n")
	f.WriteString("  font:\n")
	f.WriteString("    text: 'Roboto'\n")
	f.WriteString("    code: 'Roboto Mono'\n")
	f.WriteString("\n")
	f.WriteString("# Extensions\n")
	f.WriteString("markdown_extensions:\n")
	f.WriteString("  #- codehilite(css_class=code)\n")
	f.WriteString("  - codehilite(css_class=language-wl)\n")
	f.WriteString("  - admonition\n")

	f.Sync()
	fmt.Printf("Finished writing %v.\n", ymlFn)
}
