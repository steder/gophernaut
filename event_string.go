// Code generated by "stringer -type=Event"; DO NOT EDIT

package gophernaut

import "fmt"

const _Event_name = "StartShutdownPiningForTheFjords"

var _Event_index = [...]uint8{0, 5, 13, 31}

func (i Event) String() string {
	if i < 0 || i >= Event(len(_Event_index)-1) {
		return fmt.Sprintf("Event(%d)", i)
	}
	return _Event_name[_Event_index[i]:_Event_index[i+1]]
}
