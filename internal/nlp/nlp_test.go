package nlp

import (
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func TestNlp(t *testing.T) {
	t.Run("CapitalizeParagraphs", func(t *testing.T) {
		testcases := []struct {
			want []string
		}{
			{
				want: []string{
					"Barbara had been waiting at the table for twenty minutes. It had been twenty long and excruciating minutes. David had promised that he would be on time today. He never was, but he had promised this one time. She had made him repeat the promise multiple times over the last week until she'd believed his promise. Now she was paying the price.",
				},
			},
			{
				want: []string{
					"The red line moved across the page. With each millimeter it advanced forward, something changed in the room. The actual change taking place was difficult to perceive, but the change was real. The red line continued relentlessly across the page and the room would never be the same.",
					"She sat in the darkened room waiting. It was now a standoff. He had the power to put her in the room, but not the power to make her repent. It wasn't fair and no matter how long she had to endure the darkness, she wouldn't change her attitude. At three years old, sandy's stubborn personality had already bloomed into full view.",
				},
			},
			{
				want: []string{
					"There once lived an old man and an old woman who were peasants and had to work hard to earn their daily bread. The old man used to go to fix fences and do other odd jobs for the farmers around, and while he was gone the old woman, his wife, did the work of the house and worked in their own little plot of land.",
					"Frank knew there was a correct time and place to reveal his secret and this wasn't it. The issue was that the secret might be revealed despite his best attempt to keep it from coming out. At this point, it was out of his control and completely dependant on those around him who also knew the secret. They wouldn't purposely reveal it, or at least he believed that, but they could easily inadvertently expose it. It was going to be a long hour as he nervously eyed everyone around the table hoping they would keep their mouths shut.",
					"The wave crashed and hit the sandcastle head-on. The sandcastle began to melt under the waves force and as the wave receded, half the sandcastle was gone. The next wave hit, not quite as strong, but still managed to cover the remains of the sandcastle and take more of it away. The third wave, a big one, crashed over the sandcastle completely covering and engulfing it. When it receded, there was no trace the sandcastle ever existed and hours of hard work disappeared forever.",
				},
			},
		}
		for i, tc := range testcases {
			lower := []string{}
			for _, s := range tc.want {
				lower = append(lower, strings.ToLower(s))
			}

			actual := CapitalizeParagraphs(lower)

			if !slices.Equal(actual, tc.want) {
				t.Fatalf("tc #%d:\nwanted %#v\n\ngot %#v", i, tc.want, actual)
			}
		}
	})
}
