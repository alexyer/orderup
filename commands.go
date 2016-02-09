package main

// All bot command handlers located here.

// create-restaurant command.
// create-restaurant [restaurant name]
func createRestaurant(cmd *Cmd) string {
	switch {
	case len(cmd.Args) == 0:
		return "Restaurant name is not given."
	case len(cmd.Args) != 1:
		return "Spaces are not allowed in restaurant name."
	}

	return "created"
}
