package notes

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	notessimulation "github.com/achelabov/cosmossdk-example/x/notes/simulation"
	"github.com/achelabov/cosmossdk-example/x/notes/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	notesGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&notesGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateNote          = "op_weight_msg_notes"
		defaultWeightMsgCreateNote int = 100
	)

	var weightMsgCreateNote int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateNote, &weightMsgCreateNote, nil,
		func(_ *rand.Rand) {
			weightMsgCreateNote = defaultWeightMsgCreateNote
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateNote,
		notessimulation.SimulateMsgCreateNote(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
