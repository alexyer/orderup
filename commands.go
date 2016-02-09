package main

// All bot command handlers located here.

// create-restaurant command.
// create-restaurant [restaurant name]
func (o *Orderup) createRestaurant(cmd *Cmd) string {
	switch {
	case len(cmd.Args) == 0:
		return "Restaurant name is not given."
	case len(cmd.Args) != 1:
		return "Spaces are not allowed in restaurant name."
	}

	return "created"
}

func (o *Orderup) help(cmd *Cmd) string {
	return `Available commands:
				/orderup create-restaurant [name] -- Create a list of order numbers for restaurant name.`
}
