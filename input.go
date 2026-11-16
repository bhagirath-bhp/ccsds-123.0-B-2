package ccsds_spike

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type WeightInitMode bool

const (
	DefaultWeightInit = true
	CustomWeightInit  = false
)

type weightConfig struct {
	wInitMode          WeightInitMode `toml:"w_init_mode"`
	wCompRes           int            `toml:"w_comp_res"`
	tInc               int            `toml:"t_inc"`
	vMin               int            `toml:"v_min"`
	vMax               int            `toml:"v_max"`
	q                  int            `toml:"q"`
	expOffsetMode      bool           `toml:"exp_offset_mode"`
	expOffsetInterBVec [][]int        `toml:"exp_offset_inter_b_vec"`
	expOffsetIntraBVec []int          `toml:"exp_offset_intra_b_vec"`
	customInit         bool           `toml:"custom_init"`
	wInitVec           [][]int        `toml:"w_init_vec"`
}

type LocalSumMode int

const (
	WideNeighborOriented LocalSumMode = iota
	NarrowNeighborOriented
	WideColumnOriented
	NarrowColumnOriented
)

type PredictionMode bool

const (
	FullPredictionMode    PredictionMode = true
	ReducedPredictionMode PredictionMode = false
)

type CompressionMode bool

const (
	Lossless CompressionMode = true
	Lossy    CompressionMode = false
)

type predictionConfig struct {
	pMode      PredictionMode `toml:"p_mode"`
	sumMode    LocalSumMode   `toml:"sum_mode"`
	nPrevBands int            `toml:"n_prev_bands"`
}

type genConfig struct {
	rSize           int             `toml:"r_size"`
	imagePath       string          `toml:"image_path"`
	dRange          int             `toml:"d_range"`
	compressionMode CompressionMode `toml:"is_lossless"`
}

type srConfig struct {
	srRes     int `toml:"sr_res"`
	srDamping int `toml:"sr_damping"`
	srOffset  int `toml:"sr_offset"`
}

type errorControl struct {
	aZ int `toml:"err_lim_abs"`
	rZ int `toml:"err_lim_rel"`
}

type Config struct {
	weightConfig     weightConfig     `toml:"weights"`
	predictionConfig predictionConfig `toml:"prediction"`
	genConfig        genConfig        `toml:"general"`
	srConfig         srConfig         `toml:"sample_representatives"`
	errorControl     errorControl     `toml:"error_control"`
}

func (ccsds *CCSDS) fetchConfig(confPath string) (Config, error) {
	confBytes, err := os.ReadFile(confPath)
	if err != nil {
		return Config{}, fmt.Errorf("Error while reading config file!", err)
	}
	var config Config
	err = toml.Unmarshal(confBytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Error while parsing config", err)
	}
	return config, nil
}

func (ccsds *CCSDS) validateInput(config Config, image Image) error {
	return nil
}
