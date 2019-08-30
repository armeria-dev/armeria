package armeria

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type Command struct {
	Name         string
	AltNames     []string
	Help         string
	Hidden       bool
	Alias        string
	Permissions  *CommandPermissions
	AllowedRoles []int
	Arguments    []*CommandArgument
	Subcommands  []*Command
	Handler      func(r *CommandContext)
}

type CommandArgument struct {
	Name             string
	IncludeRemaining bool
	Optional         bool
	NoLog            bool
}

type CommandPermissions struct {
	RequireNoCharacter bool
	RequireCharacter   bool
	RequirePermission  string
}

type CommandContext struct {
	Command         *Command
	Player          *Player
	PlayerInitiated bool
	Character       *Character
	Args            map[string]string
}

// CheckPermissions returns whether or not a parent can see/use the command.
func (cmd *Command) CheckPermissions(p *Player) bool {
	if cmd.Permissions == nil {
		return true
	}

	if cmd.Permissions.RequireNoCharacter {
		if p.Character() != nil {
			return false
		}
	}

	if cmd.Permissions.RequireCharacter {
		if p.Character() == nil {
			return false
		}
	}

	if len(cmd.Permissions.RequirePermission) > 0 {
		if !p.Character().HasPermission(cmd.Permissions.RequirePermission) {
			return false
		}
	}
	return true
}

// ShowSubcommandHelp returns the list of sub-commands that the parent has access to as a string.
func (cmd *Command) ShowSubcommandHelp(p *Player, commandsEntered []string) string {
	if len(cmd.Subcommands) == 0 {
		return "There are no sub-commands available."
	}

	output := []string{
		cmd.Help,
		fmt.Sprintf("%s /%s &lt;sub-command&gt;\n",
			TextStyle("Syntax:", TextStyleBold),
			strings.Join(commandsEntered, " "),
		),
		TextStyle("Sub-commands:", TextStyleBold),
	}

	var rows []string
	for _, scmd := range cmd.Subcommands {
		if cmd.CheckPermissions(p) {
			rows = append(rows, TableRow(
				TableCell{content: TextStyle(scmd.Name, TextStyleBold)},
				TableCell{content: scmd.Help},
			))
		}
	}

	output = append(output, TextTable(rows...))

	return strings.Join(output, "\n")
}

func (cmd *Command) ShowArgumentHelp(p *Player, commandsEntered []string) string {
	if len(cmd.Arguments) == 0 {
		return "There are no command arguments."
	}

	var argumentStrings []string
	for _, arg := range cmd.Arguments {
		if arg.Optional {
			argumentStrings = append(argumentStrings, fmt.Sprintf("[%s]", arg.Name))
		} else {
			argumentStrings = append(argumentStrings, fmt.Sprintf("&lt;%s&gt;", arg.Name))
		}

	}

	output := []string{
		cmd.Help,
		fmt.Sprintf(
			"%s /%s %s\n",
			TextStyle("Syntax:", TextStyleBold),
			strings.Join(commandsEntered, " "),
			strings.Join(argumentStrings, " "),
		),
	}

	return strings.Join(output, "\n")
}

// ArgumentByName returns a CommandArgument that matches the argument's name.
func (cmd *Command) ArgumentByName(name string) *CommandArgument {
	for _, a := range cmd.Arguments {
		if strings.ToLower(a.Name) == strings.ToLower(name) {
			return a
		}
	}

	return nil
}

// LogCtx logs a parent using a command.
func (cmd *Command) LogCtx(ctx *CommandContext) {
	var args []string
	for k, v := range ctx.Args {
		a := cmd.ArgumentByName(k)
		if a == nil || !a.NoLog {
			args = append(args, fmt.Sprintf("%s=%s", k, v))
		}
	}

	c := "Anonymous"
	if ctx.Character != nil {
		c = ctx.Character.Name()
	}

	Armeria.log.Info("character executed command",
		zap.String("character", c),
		zap.String("command", ctx.Command.Name),
		zap.Strings("arguments", args),
	)
}
