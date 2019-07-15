package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

func TestVotes(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)

	tp := TestProposal
	proposal, err := keeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalID

	proposal.Status = types.StatusVotingPeriod
	keeper.SetProposal(ctx, proposal)

	// Test first vote
	keeper.AddVote(ctx, proposalID, TestAddrs[0], types.OptionAbstain)
	vote, found := keeper.GetVote(ctx, proposalID, TestAddrs[0])
	require.True(t, found)
	require.Equal(t, TestAddrs[0], vote.Voter)
	require.Equal(t, proposalID, vote.ProposalID)
	require.Equal(t, types.OptionAbstain, vote.Option)

	// Test change of vote
	keeper.AddVote(ctx, proposalID, TestAddrs[0], types.OptionYes)
	vote, found = keeper.GetVote(ctx, proposalID, TestAddrs[0])
	require.True(t, found)
	require.Equal(t, TestAddrs[0], vote.Voter)
	require.Equal(t, proposalID, vote.ProposalID)
	require.Equal(t, types.OptionYes, vote.Option)

	// Test second vote
	keeper.AddVote(ctx, proposalID, TestAddrs[1], types.OptionNoWithVeto)
	vote, found = keeper.GetVote(ctx, proposalID, TestAddrs[1])
	require.True(t, found)
	require.Equal(t, TestAddrs[1], vote.Voter)
	require.Equal(t, proposalID, vote.ProposalID)
	require.Equal(t, types.OptionNoWithVeto, vote.Option)

	// Test vote iterator
	votesIterator := keeper.GetVotesIterator(ctx, proposalID)
	require.True(t, votesIterator.Valid())
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(votesIterator.Value(), &vote)
	require.True(t, votesIterator.Valid())
	require.Equal(t, TestAddrs[0], vote.Voter)
	require.Equal(t, proposalID, vote.ProposalID)
	require.Equal(t, types.OptionYes, vote.Option)
	votesIterator.Next()
	require.True(t, votesIterator.Valid())
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(votesIterator.Value(), &vote)
	require.True(t, votesIterator.Valid())
	require.Equal(t, TestAddrs[1], vote.Voter)
	require.Equal(t, proposalID, vote.ProposalID)
	require.Equal(t, types.OptionNoWithVeto, vote.Option)
	votesIterator.Next()
	require.False(t, votesIterator.Valid())
	votesIterator.Close()
}