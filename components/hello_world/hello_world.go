package main

import (
	"fmt"

	"github.com/gotrix/gotrix"
)

var Component = helloWorld{}

type helloWorld struct{}

func (plugin *helloWorld) Include(cnf gotrix.ComponentParams) (string, error) {
	rows := ""
	for _, p := range cnf.Params() {
		rows += fmt.Sprintf("<tr><td>%s</td></tr>", p)
	}
	return fmt.Sprintf(
		`
		<div class="gtx-component hello-world">
			<h2>Hello %s!</h2>
			<p>This is hello_world.so component</p>
			<table>
				<thead>
				<tr>
					<th>Component parameter</th>
				<tr>
				</thead>
				<tbody>
				%s
				</tbody>
			</table>
		</div>`,
		cnf.Name(), rows), nil
}
