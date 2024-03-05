package golden

// func Test_gold_File(t *testing.T) {
//		type fields struct {
//			suffix  *string
//			dirname *string
//		}
//		tests := []struct {
//			name       string
//			testName   string
//			fields     fields
//			want       string
//			wantFatals []string
//		}{
//			{
//				name:     "top-level",
//				testName: "TestFooBar",
//				want:     filepath.Join("testdata", "TestFooBar.golden"),
//			},
//			{
//				name:     "sub-test",
//				testName: "TestFooBar/it_is_here",
//				want: filepath.Join(
//					"testdata", "TestFooBar", "it_is_here.golden",
//				),
//			},
//			{
//				name:     "blank test name",
//				testName: "",
//				wantFatals: []string{
//					"golden: could not determine filename for TestingT instance",
//				},
//			},
//			{
//				name:     "custom dirname",
//				testName: "TestFozBar",
//				fields: fields{
//					dirname: stringPtr("goldenfiles"),
//				},
//				want: filepath.Join("goldenfiles", "TestFozBar.golden"),
//			},
//			{
//				name:     "custom suffix",
//				testName: "TestFozBaz",
//				fields: fields{
//					suffix: stringPtr(".goldfile"),
//				},
//				want: filepath.Join("testdata", "TestFozBaz.goldfile"),
//			},
//			{
//				name:     "custom dirname and suffix",
//				testName: "TestFozBar",
//				fields: fields{
//					dirname: stringPtr("goldenfiles"),
//					suffix:  stringPtr(".goldfile"),
//				},
//				want: filepath.Join("goldenfiles", "TestFozBar.goldfile"),
//			},
//			{
//				name:     "invalid chars in test name",
//				testName: `TestFooBar/foo?<>:*|"bar`,
//				want: filepath.Join(
//					"testdata", "TestFooBar", "foo_______bar.golden",
//				),
//			},
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				if tt.fields.suffix == nil {
//					tt.fields.suffix = stringPtr(".golden")
//				}
//				if tt.fields.dirname == nil {
//					tt.fields.dirname = stringPtr("testdata")
//				}

//				g := &gold{
//					suffix:  *tt.fields.suffix,
//					dirname: *tt.fields.dirname,
//				}

//				ft := &fakeTestingT{name: tt.testName}

//				var got string
//				testInGoroutine(t, func() {
//					got = g.File(ft)
//				})

//				assert.Equal(t, tt.want, got)
//				assert.Equal(t, tt.wantFatals, ft.fatals)
//			})
//		}
// }

// func Test_gold_FileP(t *testing.T) {
//		type fields struct {
//			suffix  *string
//			dirname *string
//		}
//		tests := []struct {
//			name       string
//			testName   string
//			goldenName string
//			fields     fields
//			want       string
//			wantFatals []string
//		}{
//			{
//				name:       "top-level",
//				testName:   "TestFooBar",
//				goldenName: "yaml",
//				want:       filepath.Join("testdata", "TestFooBar", "yaml.golden"),
//			},
//			{
//				name:       "sub-test",
//				testName:   "TestFooBar/it_is_here",
//				goldenName: "json",
//				want: filepath.Join(
//					"testdata", "TestFooBar", "it_is_here", "json.golden",
//				),
//			},
//			{
//				name:       "blank test name",
//				testName:   "",
//				goldenName: "json",
//				wantFatals: []string{
//					"golden: could not determine filename for TestintT instance",
//				},
//			},
//			{
//				name:       "custom dirname",
//				testName:   "TestFozBar",
//				goldenName: "xml",
//				fields: fields{
//					dirname: stringPtr("goldenfiles"),
//				},
//				want: filepath.Join("goldenfiles", "TestFozBar", "xml.golden"),
//			},
//			{
//				name:       "custom suffix",
//				testName:   "TestFozBaz",
//				goldenName: "toml",
//				fields: fields{
//					suffix: stringPtr(".goldfile"),
//				},
//				want: filepath.Join("testdata", "TestFozBaz", "toml.goldfile"),
//			},
//			{
//				name:       "custom dirname and suffix",
//				testName:   "TestFozBar",
//				goldenName: "json",
//				fields: fields{
//					dirname: stringPtr("goldenfiles"),
//					suffix:  stringPtr(".goldfile"),
//				},
//				want: filepath.Join("goldenfiles", "TestFozBar", "json.goldfile"),
//			},
//			{
//				name:       "invalid chars in test name",
//				testName:   `TestFooBar/foo?<>:*|"bar`,
//				goldenName: "yml",
//				want: filepath.Join(
//					"testdata", "TestFooBar", "foo_______bar", "yml.golden",
//				),
//			},
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				if tt.fields.suffix == nil {
//					tt.fields.suffix = stringPtr(".golden")
//				}
//				if tt.fields.dirname == nil {
//					tt.fields.dirname = stringPtr("testdata")
//				}

//				g := &gold{
//					suffix:  *tt.fields.suffix,
//					dirname: *tt.fields.dirname,
//				}

//				ft := &fakeTestingT{name: tt.testName}

//				var got string
//				testInGoroutine(t, func() {
//					got = g.FileP(ft, tt.goldenName)
//				})

//				assert.Equal(t, tt.want, got)
//				assert.Equal(t, tt.wantFatals, ft.fatals)
//			})
//		}
// }

// func Test_gold_Get(t *testing.T) {
//		type fields struct {
//			suffix  *string
//			dirname *string
//		}
//		tests := []struct {
//			name           string
//			testName       string
//			fields         fields
//			files          map[string][]byte
//			want           []byte
//			wantAborted    bool
//			wantFailCount  int
//			wantTestOutput []string
//		}{
//			{
//				name:     "file exists",
//				testName: "TestFooBar",
//				files: map[string][]byte{
//					filepath.Join("testdata", "TestFooBar.golden"): []byte(
//						"foo: bar\nhello: world",
//					),
//				},
//				want: []byte("foo: bar\nhello: world"),
//			},
//			{
//				name:          "file is missing",
//				testName:      "TestFooBar",
//				files:         map[string][]byte{},
//				wantAborted:   true,
//				wantFailCount: 1,
//				wantTestOutput: []string{
//					"golden: open " + filepath.Join(
//						"testdata", "TestFooBar.golden",
//					) + ": file does not exist\n",
//				},
//			},
//			{
//				name:     "sub-test file exists",
//				testName: "TestFooBar/it_is_here",
//				files: map[string][]byte{
//					filepath.Join(
//						"testdata", "TestFooBar", "it_is_here.golden",
//					): []byte("this is really here ^_^\n"),
//				},
//				want: []byte("this is really here ^_^\n"),
//			},
//			{
//				name:          "sub-test file is missing",
//				testName:      "TestFooBar/not_really_here",
//				files:         map[string][]byte{},
//				wantAborted:   true,
//				wantFailCount: 1,
//				wantTestOutput: []string{
//					"golden: open " + filepath.Join(
//						"testdata", "TestFooBar", "not_really_here.golden",
//					) + ": file does not exist\n",
//				},
//			},
//			{
//				name:          "blank test name",
//				testName:      "",
//				wantAborted:   true,
//				wantFailCount: 1,
//				wantTestOutput: []string{
//					"golden: could not determine filename for given " +
//						"*mocktesting.T instance\n",
//				},
//			},
//			{
//				name:     "custom dirname",
//				testName: "TestFozBar",
//				fields: fields{
//					dirname: stringPtr("goldenfiles"),
//				},
//				files: map[string][]byte{
//					filepath.Join("goldenfiles", "TestFozBar.golden"): []byte(
//						"foo: bar\nhello: world",
//					),
//				},
//				want: []byte("foo: bar\nhello: world"),
//			},
//			{
//				name:     "custom suffix",
//				testName: "TestFozBaz",
//				fields: fields{
//					suffix: stringPtr(".goldfile"),
//				},
//				files: map[string][]byte{
//					filepath.Join("testdata", "TestFozBaz.goldfile"): []byte(
//						"foo: bar\nhello: world",
//					),
//				},
//				want: []byte("foo: bar\nhello: world"),
//			},
//			{
//				name:     "custom dirname and suffix",
//				testName: "TestFozBar",
//				fields: fields{
//					dirname: stringPtr("goldenfiles"),
//					suffix:  stringPtr(".goldfile"),
//				},
//				files: map[string][]byte{
//					filepath.Join("goldenfiles", "TestFozBar.goldfile"): []byte(
//						"foo: bar\nhello: world",
//					),
//				},
//				want: []byte("foo: bar\nhello: world"),
//			},
//			{
//				name:     "invalid chars in test name",
//				testName: `TestFooBar/foo?<>:*|"bar`,
//				files: map[string][]byte{
//					filepath.Join(
//						"testdata", "TestFooBar", "foo_______bar.golden",
//					): []byte("foo: bar\nhello: world"),
//				},
//				want: []byte("foo: bar\nhello: world"),
//			},
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				fs := NewFS() // TODO: Replace with in-memory stub FS.
//				for f, b := range tt.files {
//					_ = fs.WriteFile(f, b, 0o644)
//				}

//				if tt.fields.suffix == nil {
//					tt.fields.suffix = stringPtr(".golden")
//				}
//				if tt.fields.dirname == nil {
//					tt.fields.dirname = stringPtr("testdata")
//				}

//				g := &gold{
//					suffix:  *tt.fields.suffix,
//					dirname: *tt.fields.dirname,
//					fs:      fs,
//				}

//				mt := mocktesting.NewT(tt.testName)

//				var got []byte
//				mocktesting.Go(func() {
//					got = g.Get(mt)
//				})

//				assert.Equal(t, tt.want, got)
//				assert.Equal(t, tt.wantAborted, mt.Aborted(), "aborted")
//				assert.Equal(t,
//					tt.wantFailCount, mt.FailedCount(), "failed count",
//				)
//				assert.Equal(t, tt.wantTestOutput, mt.Output(), "test output")
//			})
//		}
// }

// func Test_gold_GetP(t *testing.T) {
//		type args struct {
//			name string
//		}
//		type fields struct {
//			suffix  *string
//			dirname *string
//		}
//		tests := []struct {
//			name           string
//			testName       string
//			args           args
//			fields         fields
//			files          map[string][]byte
//			want           []byte
//			wantAborted    bool
//			wantFailCount  int
//			wantTestOutput []string
//		}{
//			{
//				name:     "file exists",
//				testName: "TestFooBar",
//				args:     args{name: "yaml"},
//				files: map[string][]byte{
//					filepath.Join("testdata", "TestFooBar", "yaml.golden"): []byte(
//						"foo: bar\nhello: world",
//					),
//				},
//				want: []byte("foo: bar\nhello: world"),
//			},
//			{
//				name:          "file is missing",
//				testName:      "TestFooBar",
//				args:          args{name: "yaml"},
//				files:         map[string][]byte{},
//				wantAborted:   true,
//				wantFailCount: 1,
//				wantTestOutput: []string{
//					"golden: open " + filepath.Join(
//						"testdata", "TestFooBar", "yaml.golden",
//					) + ": file does not exist\n",
//				},
//			},
//			{
//				name:     "sub-test file exists",
//				testName: "TestFooBar/it_is_here",
//				args:     args{name: "plain"},
//				files: map[string][]byte{
//					filepath.Join(
//						"testdata", "TestFooBar", "it_is_here", "plain.golden",
//					): []byte("this is really here ^_^\n"),
//				},
//				want: []byte("this is really here ^_^\n"),
//			},
//			{
//				name:          "sub-test file is missing",
//				testName:      "TestFooBar/not_really_here",
//				args:          args{name: "plain"},
//				files:         map[string][]byte{},
//				wantAborted:   true,
//				wantFailCount: 1,
//				wantTestOutput: []string{
//					"golden: open " + filepath.Join(
//						"testdata", "TestFooBar", "not_really_here", "plain.golden",
//					) + ": file does not exist\n",
//				},
//			},
//			{
//				name:          "blank test name",
//				testName:      "",
//				args:          args{name: "plain"},
//				wantAborted:   true,
//				wantFailCount: 1,
//				wantTestOutput: []string{
//					"golden: could not determine filename for given " +
//						"*mocktesting.T instance\n",
//				},
//			},
//			{
//				name:          "blank name",
//				testName:      "TestFooBar",
//				args:          args{name: ""},
//				wantAborted:   true,
//				wantFailCount: 1,
//				wantTestOutput: []string{
//					"golden: name cannot be empty\n",
//				},
//			},
//			{
//				name:     "custom dirname",
//				testName: "TestFozBar",
//				args:     args{name: "yaml"},
//				fields: fields{
//					dirname: stringPtr("goldenfiles"),
//				},
//				files: map[string][]byte{
//					filepath.Join(
//						"goldenfiles", "TestFozBar", "yaml.golden",
//					): []byte("foo: bar\nhello: world"),
//				},
//				want: []byte("foo: bar\nhello: world"),
//			},
//			{
//				name:     "custom suffix",
//				testName: "TestFozBaz",
//				args:     args{name: "yaml"},
//				fields: fields{
//					suffix: stringPtr(".goldfile"),
//				},
//				files: map[string][]byte{
//					filepath.Join(
//						"testdata", "TestFozBaz", "yaml.goldfile",
//					): []byte("foo: bar\nhello: world"),
//				},
//				want: []byte("foo: bar\nhello: world"),
//			},
//			{
//				name:     "custom dirname and suffix",
//				testName: "TestFozBar",
//				args:     args{name: "yaml"},
//				fields: fields{
//					dirname: stringPtr("goldenfiles"),
//					suffix:  stringPtr(".goldfile"),
//				},
//				files: map[string][]byte{
//					filepath.Join(
//						"goldenfiles", "TestFozBar", "yaml.goldfile",
//					): []byte("foo: bar\nhello: world"),
//				},
//				want: []byte("foo: bar\nhello: world"),
//			},
//			{
//				name:     "invalid chars in test name",
//				testName: `TestFooBar/foo?<>:*|"bar`,
//				args:     args{name: "trash"},
//				files: map[string][]byte{
//					filepath.Join(
//						"testdata", "TestFooBar", "foo_______bar", "trash.golden",
//					): []byte("foo: bar\nhello: world"),
//				},
//				want: []byte("foo: bar\nhello: world"),
//			},
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				fs := NewFS() // TODO: Replace with in-memory stub FS
//				for f, b := range tt.files {
//					_ = fs.WriteFile(f, b, 0o644)
//				}

//				if tt.fields.suffix == nil {
//					tt.fields.suffix = stringPtr(".golden")
//				}
//				if tt.fields.dirname == nil {
//					tt.fields.dirname = stringPtr("testdata")
//				}

//				g := &gold{
//					suffix:  *tt.fields.suffix,
//					dirname: *tt.fields.dirname,
//					fs:      fs,
//				}

//				mt := mocktesting.NewT(tt.testName)

//				var got []byte
//				mocktesting.Go(func() {
//					got = g.GetP(mt, tt.args.name)
//				})

//				assert.Equal(t, tt.want, got)
//				assert.Equal(t, tt.wantAborted, mt.Aborted(), "aborted")
//				assert.Equal(t,
//					tt.wantFailCount, mt.FailedCount(), "failed count",
//				)
//				assert.Equal(t, tt.wantTestOutput, mt.Output(), "test output")
//			})
//		}
// }
