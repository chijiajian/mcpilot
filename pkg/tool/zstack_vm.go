package tool

import "fmt"

type ZStackVMTool struct{}

func (t *ZStackVMTool) Name() string {
	return "ZStackVMTool"
}

func (t *ZStackVMTool) Description() string {
	return "A tool for managing ZStack virtual machines."
}

func (t *ZStackVMTool) InputSchema() map[string]string {
	return map[string]string{
		"name":   "虚拟机名称（字符串）",
		"cpu":    "CPU 核心数（整数）",
		"memory": "内存大小（单位：MB，整数）",
	}
}

func (t *ZStackVMTool) Run(params map[string]string) (string, error) {
	name := params["name"]
	cpu := params["cpu"]
	memory := params["memory"]

	// TODO: Implement the zstack terraform zstack provider to create a VM
	return fmt.Sprintf("VM '%s' Success deployed [CPU: %s, Mem: %sMB]", name, cpu, memory), nil
}
