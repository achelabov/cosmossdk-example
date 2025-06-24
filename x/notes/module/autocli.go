package notes

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/achelabov/cosmossdk-example/x/notes/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListNote",
					Use:       "list-note",
					Short:     "List all note",
				},
				{
					RpcMethod:      "GetNote",
					Use:            "get-note [id]",
					Short:          "Gets a note",
					Alias:          []string{"show-note"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				{
					RpcMethod:      "NotesByCreator",
					Use:            "notes-by-creator [creator]",
					Short:          "Query notes-by-creator",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "creator"}},
				},

				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateNote",
					Use:            "create-note [text]",
					Short:          "Send a create-note tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "text"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
