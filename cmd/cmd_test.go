package cmd

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/noborus/trdsql"
)

func Test_inputFormat(t *testing.T) {
	type args struct {
		i inputFlag
	}
	tests := []struct {
		name string
		args args
		want trdsql.Format
	}{
		{
			name: "testCSV",
			args: args{
				i: inputFlag{
					CSV: true,
				},
			},
			want: trdsql.CSV,
		},
		{
			name: "testLTSV",
			args: args{
				i: inputFlag{
					LTSV: true,
					JSON: true,
				},
			},
			want: trdsql.LTSV,
		},
		{
			name: "testJSON",
			args: args{
				i: inputFlag{
					JSON: true,
					TBLN: true,
				},
			},
			want: trdsql.JSON,
		},
		{
			name: "testTBLN",
			args: args{
				i: inputFlag{
					TBLN: true,
				},
			},
			want: trdsql.TBLN,
		},
		{
			name: "testGUESS",
			args: args{
				i: inputFlag{},
			},
			want: trdsql.GUESS,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inputFormat(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inputFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_outputFormat(t *testing.T) {
	type args struct {
		o outputFlag
	}
	tests := []struct {
		name string
		args args
		want trdsql.Format
	}{
		{
			name: "testCSV",
			args: args{
				o: outputFlag{
					CSV: true,
				},
			},
			want: trdsql.CSV,
		},
		{
			name: "testLTSV",
			args: args{
				o: outputFlag{
					CSV:  false,
					LTSV: true,
				},
			},
			want: trdsql.LTSV,
		},
		{
			name: "testAT",
			args: args{
				o: outputFlag{
					CSV:  false,
					LTSV: false,
					AT:   true,
				},
			},
			want: trdsql.AT,
		},
		{
			name: "testMD",
			args: args{
				o: outputFlag{
					CSV:  false,
					LTSV: false,
					MD:   true,
				},
			},
			want: trdsql.MD,
		},
		{
			name: "testVF",
			args: args{
				o: outputFlag{
					CSV:  false,
					LTSV: false,
					VF:   true,
				},
			},
			want: trdsql.VF,
		},
		{
			name: "testRAW",
			args: args{
				o: outputFlag{
					CSV:  false,
					LTSV: false,
					RAW:  true,
				},
			},
			want: trdsql.RAW,
		},
		{
			name: "testJSON",
			args: args{
				o: outputFlag{
					CSV:  false,
					LTSV: false,
					JSON: true,
				},
			},
			want: trdsql.JSON,
		},
		{
			name: "testJSONL",
			args: args{
				o: outputFlag{
					CSV:   false,
					LTSV:  false,
					JSONL: true,
				},
			},
			want: trdsql.JSONL,
		},
		{
			name: "testTBLN",
			args: args{
				o: outputFlag{
					TBLN: true,
				},
			},
			want: trdsql.TBLN,
		},
		{
			name: "testDEFAULT",
			args: args{
				o: outputFlag{},
			},
			want: trdsql.CSV,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := outputFormat(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("outputFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getQuery(t *testing.T) {
	type argss struct {
		args     []string
		fileName string
	}
	tests := []struct {
		name    string
		argss   argss
		want    string
		wantErr bool
	}{
		{
			name: "testARGS",
			argss: argss{
				[]string{"SELECT 1"},
				"",
			},
			want:    "SELECT 1",
			wantErr: false,
		},
		{
			name: "testARGS2",
			argss: argss{
				[]string{"SELECT", "1"},
				"",
			},
			want:    "SELECT 1",
			wantErr: false,
		},
		{
			name: "testTrim",
			argss: argss{
				[]string{"SELECT * FROM test;   "},
				"",
			},
			want:    "SELECT * FROM test",
			wantErr: false,
		},
		{
			name: "testFileErr",
			argss: argss{
				[]string{},
				filepath.Join("..", "testdata", "noFile.sql"),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "testFile",
			argss: argss{
				[]string{},
				filepath.Join("..", "testdata", "test.sql"),
			},
			want:    "SELECT * FROM testdata/test.csv",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getQuery(tt.argss.args, tt.argss.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDB(t *testing.T) {
	type argss struct {
		cfg     *config
		cDB     string
		cDriver string
		cDSN    string
	}
	tests := []struct {
		name  string
		argss argss
		want  string
		want1 string
	}{
		{
			name: "testNoConfig",
			argss: argss{
				cfg:     &config{},
				cDB:     "",
				cDriver: "postgres",
				cDSN:    "dbname=test",
			},
			want:  "postgres",
			want1: "dbname=test",
		},
		{
			name: "testNoConfigDB",
			argss: argss{
				cfg:     &config{},
				cDB:     "test",
				cDriver: "postgres",
				cDSN:    "dbname=\"test\"",
			},
			want:  "postgres",
			want1: "dbname=\"test\"",
		},
		{
			name: "testDSN",
			argss: argss{
				cfg:     &config{},
				cDB:     "",
				cDriver: "",
				cDSN:    "dbname=\"test\"",
			},
			want:  "",
			want1: "dbname=\"test\"",
		},
		{
			name: "testConfig",
			argss: argss{
				cfg: &config{
					Db: "",
					Database: map[string]database{
						"pdb": {
							Driver: "postgres",
							Dsn:    "dbname=\"test\"",
						},
					},
				},
				cDB:     "pdb",
				cDriver: "",
				cDSN:    "",
			},
			want:  "postgres",
			want1: "dbname=\"test\"",
		},
		{
			name: "testConfigErr",
			argss: argss{
				cfg: &config{
					Db: "",
					Database: map[string]database{
						"pdb": {
							Driver: "postgres",
							Dsn:    "dbname=\"test\"",
						},
					},
				},
				cDB:     "sdb",
				cDriver: "",
				cDSN:    "",
			},
			want:  "",
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			got, got1 := getDB(tt.argss.cfg, tt.argss.cDB, tt.argss.cDriver, tt.argss.cDSN)
			if got != tt.want {
				t.Errorf("getDB() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getDB() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_optsCommand(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "testEmpty",
			args: []string{"trdsql", "-a"},
			want: "trdsql",
		},
		{
			name: "testFile",
			args: []string{"trdsql", "-ih", "-a", "test.csv"},
			want: "trdsql -ih",
		},
		{
			name: "testFile2",
			args: []string{"trdsql", "-ih", "-ir", "2", "-a", "test.csv"},
			want: "trdsql -ih -ir 2",
		},
		{
			name: "testStdin",
			args: []string{"trdsql", "-ih", "-a", "-"},
			want: "trdsql -ih",
		},
		{
			name: "testFile",
			args: []string{"trdsql", "-dsn=\"dbname=test\"", "-a", "test.csv"},
			want: "trdsql -dsn=\"dbname=test\"",
		},
		{
			name: "testDelimiterSpace",
			args: []string{"trdsql", "-id", " ", "-a", "test.csv"},
			want: "trdsql -id \" \"",
		},
		{
			name: "testDelimiterUnder",
			args: []string{"trdsql", "-id", "_", "-a", "test.csv"},
			want: "trdsql -id _",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := trdsql.NewAnalyzeOpts()
			got := optsCommand(opts, tt.args)
			if !reflect.DeepEqual(got.Command, tt.want) {
				t.Errorf("optsCommand() = %v, want %v", got.Command, tt.want)
			}
		})
	}
}

func TestCli_Run(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want int
	}{
		{
			name: "testEmpty",
			args: []string{"trdsql"},
			want: 2,
		},
		{
			name: "testOne",
			args: []string{"trdsql", "SELECT 1"},
			want: 0,
		},
		{
			name: "testErr",
			args: []string{"trdsql", "Err"},
			want: 1,
		},
		{
			name: "testAnalyze",
			args: []string{"trdsql", "-a", filepath.Join("..", "testdata", "test.csv")},
			want: 0,
		},
		{
			name: "testAnalyze2",
			args: []string{"trdsql", "-ir", "1", "-a", filepath.Join("..", "testdata", "test.csv")},
			want: 0,
		},
		{
			name: "testSQLOnly",
			args: []string{"trdsql", "-A", filepath.Join("..", "testdata", "test.csv")},
			want: 0,
		},
		{
			name: "testDebug",
			args: []string{"trdsql", "-debug", "SELECT 1"},
			want: 0,
		},
		{
			name: "testDBList",
			args: []string{"trdsql", "-dblist"},
			want: 0,
		},
		{
			name: "testHelp",
			args: []string{"trdsql", "-help"},
			want: 2,
		},
		{
			name: "testVersion",
			args: []string{"trdsql", "-version"},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
			cli := Cli{
				OutStream: outStream,
				ErrStream: errStream,
			}
			var buf bytes.Buffer
			log.SetOutput(&buf)
			if got := cli.Run(tt.args); got != tt.want {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_printDBList(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config
	}{
		{
			name: "test",
			cfg: &config{
				Db: "",
				Database: map[string]database{
					"pdb": {Driver: "postgres", Dsn: "dbname=test"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printDBList(tt.cfg)
		})
	}
}

func Test_colorOpts(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping specific test")
	}
	tests := []struct {
		name   string
		setEnv bool
		want   bool
	}{
		{
			name:   "test",
			setEnv: false,
			want:   true,
		},
		{
			name:   "test",
			setEnv: true,
			want:   false,
		},
	}
	for _, tt := range tests {
		opts := trdsql.NewAnalyzeOpts()
		if tt.setEnv {
			os.Setenv("NO_COLOR", "1")
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := colorOpts(opts); !reflect.DeepEqual(got.Color, tt.want) {
				t.Errorf("colorOpts() = %v, want %v", got.Color, tt.want)
			}
		})
	}
}

func Test_quoteOpts(t *testing.T) {
	tests := []struct {
		name   string
		driver string
		want   string
	}{
		{
			name:   "testSQLIte3",
			driver: "sqlite3",
			want:   "\\`",
		},
		{
			name:   "testMySQL",
			driver: "mysql",
			want:   "\\`",
		},
		{
			name:   "testPostgreSQL",
			driver: "postgres",
			want:   `\"`,
		},
	}
	for _, tt := range tests {
		opts := trdsql.NewAnalyzeOpts()
		t.Run(tt.name, func(t *testing.T) {
			got := quoteOpts(opts, tt.driver)
			if !reflect.DeepEqual(got.Quote, tt.want) {
				t.Errorf("quoteOpts() = %v, want %v", got.Quote, tt.want)
			}
		})
	}
}
