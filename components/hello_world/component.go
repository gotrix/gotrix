package component

import (
	"fmt"

	"github.com/gotrix/gotrix"
)

type Component struct{}

func (*Component) Include(cnf gotrix.ComponentParams) (string, error) {
	rows := ""
	for _, p := range cnf.Params() {
		rows += fmt.Sprintf("<tr><td>%s</td></tr>", p)
	}
	return fmt.Sprintf(
		`
		<div class="hello-world">
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
