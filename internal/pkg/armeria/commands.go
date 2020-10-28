package armeria

import (
	"armeria/internal/pkg/misc"
	"strings"
)

const (
	CommandErrNoPerms = "You cannot use that command."
	CommandErrInvalid = "That's an invalid command."
)

// Manager is the global manager instance for Command objects
type CommandManager struct {
	commands []*Command
}

// NewCommandManager will return a new instance of the command manager.
func NewCommandManager() *CommandManager {
	return &CommandManager{
		commands: []*Command{},
	}
}

// Commands returns all the registered commands in the game.
func (m *CommandManager) Commands() []*Command {
	return m.commands
}

// RegisterCommand will register a Command with the command manager with the arguments
// parsed out.
func (m *CommandManager) RegisterCommand(c *Command) {
	// Set parents for sub-commands.
	for _, cmd := range c.Subcommands {
		cmd.Parent = c
	}

	m.commands = append(m.commands, c)
}

// FindCommand will return a matched registered Command.
func (m *CommandManager) FindCommand(p *Player, searchWithin []*Command, cmd string, alreadyProcessed []string) (*Command, map[string]string, string) {
	sections := strings.Fields(cmd)
	cmdName := strings.ToLower(sections[0])

	for _, cmd := range searchWithin {
		if strings.ToLower(cmd.Name) == cmdName || misc.Contains(cmd.AltNames, cmdName) {
			// Handle permissions
			if !cmd.CheckPermissions(p) {
				return nil, nil, CommandErrNoPerms
			}

			// Handle sub-commands
			if cmd.Subcommands != nil {
				processedCommands := append(alreadyProcessed, cmdName)
				if len(sections) == 1 {
					return nil, nil, cmd.ShowSubcommandHelp(p, processedCommands)

				}
				return m.FindCommand(p, cmd.Subcommands, strings.Join(sections[1:], " "), processedCommands)
			}

			// Parse and store arguments, if any
			commandArgs := make(map[string]string)
			parsedArgs := misc.ParseArguments(sections[1:])
			if cmd.Arguments != nil {
				if len(parsedArgs) > 0 && parsedArgs[len(parsedArgs)-1] == "--help" {
					return nil, nil, cmd.ShowArgumentHelp(append(alreadyProcessed, cmdName))
				}

				for pos, arg := range cmd.Arguments {
					if !arg.Optional && len(parsedArgs) < (pos+1) {
						return nil, nil, cmd.ShowArgumentHelp(append(alreadyProcessed, cmdName))
					}
					if arg.IncludeRemaining {
						commandArgs[arg.Name] = strings.Join(parsedArgs[pos:], " ")
					} else if len(parsedArgs) >= pos+1 {
						commandArgs[arg.Name] = parsedArgs[pos]
					} else {
						commandArgs[arg.Name] = ""
					}
				}
			}

			return cmd, commandArgs, ""
		}
	}

	return nil, nil, CommandErrInvalid
}

// ProcessCommand will evaluate and process a command sent by the parent either
// manually or programmatically.
func (m *CommandManager) ProcessCommand(p *Player, command string, playerInitiated bool) {
	sections := strings.Fields(command)
	if len(sections) == 0 {
		return
	}

	// short-circuit for channel commands
	for _, ch := range Armeria.channels {
		if ch.SlashCommand == "/"+strings.ToLower(sections[0]) {
			if ch.HasPermission(p.Character()) {
				if p.Character().InChannel(ch) {
					remainingText := sections[1:]
					sections = []string{"channel", "say", ch.Name}
					sections = append(sections, remainingText...)
				} else {
					p.client.ShowColorizedText("You are not in that channel.", ColorError)
					return
				}
			} else {
				p.client.ShowColorizedText("You are not able to chat on that channel.", ColorError)
				return
			}
		}
	}
	cmd, cmdArgs, errorMsg := m.FindCommand(p, m.commands, strings.Join(sections, " "), []string{})

	if cmd == nil {
		p.client.ShowColorizedText(errorMsg, ColorCmdHelp)
		return
	}

	ctx := &CommandContext{
		Command:         cmd,
		Player:          p,
		Args:            cmdArgs,
		PlayerInitiated: playerInitiated,
	}

	if p.Character() != nil {
		ctx.Character = p.Character()
	}

	if len(cmd.Alias) > 0 {
		m.ProcessCommand(p, cmd.Alias, playerInitiated)
		return
	}

	cmd.LogCtx(ctx)
	cmd.Handler(ctx)
}
