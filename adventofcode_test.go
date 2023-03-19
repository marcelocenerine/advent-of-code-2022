package adventofcode

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSolutions(t *testing.T) {
	tests := []struct {
		puzzle Puzzle
		want   Result
	}{
		{
			puzzle: CalorieCounting{},
			want:   Result{Part1: "72511", Part2: "212117"},
		},
		{
			puzzle: RockPaperScissors{},
			want:   Result{Part1: "12679", Part2: "14470"},
		},
		{
			puzzle: RucksackReorganization{},
			want:   Result{Part1: "7737", Part2: "2697"},
		},
		{
			puzzle: CampCleanup{},
			want:   Result{Part1: "475", Part2: "825"},
		},
		{
			puzzle: SupplyStacks{},
			want:   Result{Part1: "TLFGBZHCN", Part2: "QRQFHFWCL"},
		},
		{
			puzzle: TuningTrouble{},
			want:   Result{Part1: "1142", Part2: "2803"},
		},
		{
			puzzle: NoSpaceLeftOnDevice{},
			want:   Result{Part1: "1443806", Part2: "942298"},
		},
		{
			puzzle: TreetopTreeHouse{},
			want:   Result{Part1: "1700", Part2: "470596"},
		},
		{
			puzzle: RopeBridge{},
			want:   Result{Part1: "6498", Part2: "2531"},
		},
		{
			puzzle: CathodeRayTube{},
			want: Result{
				Part1: "13680",
				Part2: `###..####..##..###..#..#.###..####.###..
#..#....#.#..#.#..#.#.#..#..#.#....#..#.
#..#...#..#....#..#.##...#..#.###..###..
###...#...#.##.###..#.#..###..#....#..#.
#....#....#..#.#....#.#..#....#....#..#.
#....####..###.#....#..#.#....####.###..`,
			},
		},
		{
			puzzle: MonkeyInTheMiddle{},
			want:   Result{Part1: "98280", Part2: "17673687232"},
		},
		{
			puzzle: HillClimbingAlgorithm{},
			want:   Result{Part1: "?", Part2: "?"},
		},
		{
			puzzle: DistressSignal{},
			want:   Result{Part1: "?", Part2: "?"},
		},
		{
			puzzle: RegolithReservoir{},
			want:   Result{Part1: "?", Part2: "?"},
		},
		{
			puzzle: BeaconExclusionZone{},
			want:   Result{Part1: "?", Part2: "?"},
		},
		{
			puzzle: ProboscideaVolcanium{},
			want:   Result{Part1: "?", Part2: "?"},
		},
		{
			puzzle: PyroclasticFlow{},
			want:   Result{Part1: "?", Part2: "?"},
		},
		{
			puzzle: BoilingBoulders{},
			want:   Result{Part1: "?", Part2: "?"},
		},
		{
			puzzle: NotEnoughMinerals{},
			want:   Result{Part1: "?", Part2: "?"},
		},
		{
			puzzle: GrovePositioningSystem{},
			want:   Result{Part1: "?", Part2: "?"},
		},
		{
			puzzle: MonkeyMath{},
			want:   Result{Part1: "?", Part2: "?"},
		},
		{
			puzzle: MonkeyMap{},
			want:   Result{Part1: "?", Part2: "?"},
		},
		{
			puzzle: UnstableDiffusion{},
			want:   Result{Part1: "?", Part2: "?"},
		},
		{
			puzzle: BlizzardBasin{},
			want:   Result{Part1: "?", Part2: "?"},
		},
		{
			puzzle: FullOfHotAir{},
			want:   Result{Part1: "?", Part2: "?"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.puzzle.Details().String(), func(t *testing.T) {
			input, err := LoadInput(tc.puzzle)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got, err := tc.puzzle.Solve(&input)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}
