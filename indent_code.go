// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lute

func (t *Tree) parseIndentCode(line items) (ret Node) {
	baseNode := &BaseNode{typ: NodeCode}
	code := &Code{baseNode, "", ""}
	var codeValue string
Loop:
	for {
		var spaces, tabs int
		for i := 0; i < 4; i++ {
			token := line[i]
			if itemSpace == token.typ {
				spaces++
			} else if itemTab == token.typ {
				tabs++
			}
			if 3 < spaces || 0 < tabs {
				line = line[i+1:]
				break
			}
		}

		if line.isBlankLine() {
			return
		}

		for i := 0; i < len(line); i++ {
			token := line[i]
			codeValue += token.val
			if token.isNewline() {
				newlines, nonNewline := t.nonNewline()
				if nonNewline.isEOF() {
					break Loop
				}

				if t.blockquoteMarkerCount(line) < t.context.BlockquoteLevel {
					t.backupLine(nonNewline)
					break
				}

				line = nonNewline
				spaces, tabs, _, _ := t.nonWhitespace(line)
				if 1 > tabs && 4 > spaces {
					t.backupLine(line)
					break Loop
				} else {
					codeValue += newlines.rawText()
					newlines = append(newlines, token)
					code.tokens = append(code.tokens, newlines...)

					continue Loop
				}
			}
			code.tokens = append(code.tokens, token)
		}

		line = t.nextLine()
		if !t.isIndentCode(line) {
			t.backupLine(line)
			break
		}

		if t.blockquoteMarkerCount(line) < t.context.BlockquoteLevel {
			t.backupLine(line)
			break
		}
	}

	code.Value = codeValue
	ret = code

	return
}

func (t *Tree) isIndentCode(line items) bool {
	var spaces int
	for _, token := range line {
		if itemSpace == token.typ {
			spaces++
			continue
		}
		if itemTab == token.typ {
			spaces += 4
			continue
		}

		break
	}

	return t.context.IndentSpaces+3 < spaces
}

func (t *Tree) nonNewline() (newlines items, line items) {
	for line = t.nextLine(); line.isBlankLine() && !line.isEOF(); line = t.nextLine() {
		if 5 > len(line) {
			_, line = line.trimLeftSpace()
		} else {
			line = t.indentOffset(line, 4)
		}

		newlines = append(newlines, line...)
	}

	return
}
