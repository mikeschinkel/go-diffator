package diffator_test

import (
	"testing"

	"github.com/mikeschinkel/go-diffator"
)

func TestCompareStrings(t *testing.T) {
	type args struct {
		s1     string
		s2     string
		pad    *diffator.IntValue
		format *diffator.StringValue
		minLen *diffator.IntValue
	}
	var tests = []struct {
		name string
		args args
		want string
	}{
		{
			name: "S1 and S2 are empty",
		},
		{
			name: "S1 is empty",
			args: args{
				s2: "ABC",
			},
			want: "<(/ABC)>",
		},
		{
			name: "S2 is empty",
			args: args{
				s1: "ABC",
			},
			want: "<(ABC/)>",
		},
		{
			name: "S1 and S2 are completely different",
			args: args{
				s1: "ABC",
				s2: "XYZ",
			},
			want: "<(ABC/XYZ)>",
		},
		{
			name: "S1 and S2 start the same, but end different",
			args: args{
				s1: "ABCDEF",
				s2: "ABCDXYZ",
			},
			want: "ABCD<(EF/XYZ)>",
		},
		{
			name: "S1 and S2 start different but end the same",
			args: args{
				s1: "ABCDXYZ",
				s2: "123XYZ",
			},
			want: "<(ABCD/123)>XYZ",
		},
		{
			name: "S1 and S2 start different but end the same",
			args: args{
				s1: "Look, it's Batman!!!",
				s2: "Look, it's Superman!!!",
			},
			want: "Look, it's <(Bat/Super)>man!!!",
		},
		{
			name: "S1 has extra middle chars",
			args: args{
				s1:  "ABCDEF123GHIJKLMNOP",
				s2:  "ABCDEFGHIJKLMNOP",
				pad: diffator.Int(5),
			},
			want: "BCDEF<(123/)>GHIJK",
		},
		{
			name: "S1 has prefix and suffix that S2 does not have",
			args: args{
				s1:     "123GHI456",
				s2:     "GHI",
				minLen: diffator.Int(2),
			},
			want: "<(123/)>GHI<(456/)>",
		},
		{
			name: "S1 and S2 share a middle, differ on the ends",
			args: args{
				s1:     "123GHI789",
				s2:     "987GHI321",
				minLen: diffator.Int(2),
			},
			want: "<(123/987)>GHI<(789/321)>",
		},
		{
			name: "S1 has two sets of extra middle chars",
			args: args{
				s1:     "ABCDEF123GHI456JKLMNOP",
				s2:     "ABCDEFGHIJKLMNOP",
				pad:    diffator.Int(5),
				minLen: diffator.Int(2),
			},
			want: "BCDEF<(123/)>GHI<(456/)>JKLMN",
		},
		{
			name: "And vs. &",
			args: args{
				s1:  "Publishing and graphic design.",
				s2:  "Publishing & graphic design.",
				pad: diffator.Int(25),
			},
			want: "Publishing <(and/&)> graphic design.",
		},
		{
			name: "Short Lorem Ipsum with format",
			args: args{
				s1:     "Lorem ipsum may be used as a placeholder before final copy is available.",
				s2:     "Lorem ipsum is often used as a placeholder awaiting final copy.",
				pad:    diffator.Int(25),
				minLen: diffator.Int(3),
				format: diffator.String("{%s|%s}"),
			},
			want: "Lorem ipsum {may be|is often} used as a placeholder {before|awaiting} final copy{ is available|}.",
		},
		{
			name: "Sans",
			args: args{
				s1:     "typeface without relying on meaningful content.",
				s2:     "typeface sans meaningful content.",
				minLen: diffator.Int(3),
			},
			want: "typeface <(without relying on/sans)> meaningful content.",
		},
		{
			name: "Longer Lorem Ipsum",
			args: args{
				s1:     "In publishing and graphic design, Lorem ipsum is a placeholder text commonly used to demonstrate the visual form of a document or a typeface without relying on meaningful content. Lorem ipsum may be used as a placeholder before final copy is available.",
				s2:     "In publishing & graphic design, Lorem ipsum is a commonly used text placeholder to demonstrate a document in its visual form, or a typeface sans meaningful content. Lorem ipsum is often used as a placeholder awaiting final copy.",
				pad:    diffator.Int(25),
				minLen: diffator.Int(3),
			},
			want: "In publishing <(and/&)> graphic design, Lorem ipsum is a <(placeholder text /)>commonly used<(/ text placeholder)> to demonstrate <(the/a document in its)> visual form<( of a document/,)> or a typeface <(without relying on/sans)> meaningful content. Lorem ipsum <(may be/is often)> used as a placeholder <(before/awaiting)> final copy<( is available/)>.",
		},
		{
			name: "Reordered substrings",
			args: args{
				s1:     "version tag already exists [project='golang'] [version_tag='go1.21.4']",
				s2:     "version tag already exists [version_tag='go1.21.4'] [project='golang']",
				minLen: diffator.Int(3),
			},
			want: "version tag already exists [<(project='golang'] [/)>version_tag='go1.21.4<(/'] [project='golang)>']",
		},
		{
			name: "Mixing word order",
			args: args{
				s1:     "Lorem ipsum is a placeholder text commonly used to demonstrate the visual form of a document.",
				s2:     "Lorem ipsum is a commonly used text placeholder to demonstrate a document in its visual form.",
				minLen: diffator.Int(3),
			},
			want: "Lorem ipsum is a <(placeholder text /)>commonly used<(/ text placeholder)> to demonstrate <(the/a document in its)> visual form<( of a document/)>.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := diffator.CompareStrings(tt.args.s1, tt.args.s2, &diffator.StringOpts{
				MatchingPadLen:  tt.args.pad,
				MinSubstrLen:    tt.args.minLen,
				LeftRightFormat: tt.args.format,
			})
			if got != tt.want {
				t.Errorf("\ndiff.CompareStrings(s1,s2):\n\t got: %v\n\twant: %v\n", got, tt.want)
			}
		})
	}
}
