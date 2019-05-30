package armeria

import (
	"fmt"
	"io/ioutil"

	"go.uber.org/zap"

	lua "github.com/yuin/gopher-lua"
)

// ReadMobScript returns the script contents for a mob from disk.
func ReadMobScript(m *Mob) string {
	b, err := ioutil.ReadFile(m.GetScriptFile())
	if err != nil {
		return ""
	}
	return string(b)
}

// WriteMobScript writes a mob script to disk.
func WriteMobScript(m *Mob, script string) {
	_ = ioutil.WriteFile(m.GetScriptFile(), []byte(script), 0644)
}

func LuaMobSay(L *lua.LState) int {
	text := L.ToString(1)
	mname := lua.LVAsString(L.GetGlobal("mob_name"))
	mid := lua.LVAsString(L.GetGlobal("mob_instance"))

	m := Armeria.mobManager.GetMobByName(mname)
	mi := m.GetInstanceById(mid)

	for _, c := range mi.GetRoom().GetCharacters(nil) {
		c.GetPlayer().clientActions.ShowColorizedText(
			fmt.Sprintf("%s says, \"%s\".", mi.GetFName(), text),
			ColorSay,
		)
	}

	return 0
}

func CallMobFunc(invoker *Character, mi *MobInstance, funcName string) {
	L := lua.NewState()
	defer L.Close()

	// global variables
	L.SetGlobal("invoker_name", lua.LString(invoker.GetName()))
	L.SetGlobal("mob_instance", lua.LString(mi.Id))
	L.SetGlobal("mob_name", lua.LString(mi.Parent))
	// global functions
	L.SetGlobal("mob_say", L.NewFunction(LuaMobSay))

	err := L.DoFile(mi.GetParent().GetScriptFile())
	if err != nil {
		Armeria.log.Error("error compiling lua script",
			zap.String("script", mi.GetParent().GetScriptFile()),
			zap.Error(err),
		)
		return
	}

	err = L.CallByParam(lua.P{
		Fn:      L.GetGlobal(funcName),
		NRet:    0,
		Protect: true,
	})
	if err != nil {
		Armeria.log.Error("error executing function in lua script",
			zap.String("script", mi.GetParent().GetScriptFile()),
			zap.String("function", funcName),
			zap.Error(err),
		)
	}
}