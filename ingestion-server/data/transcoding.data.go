package data

import (
	"os/exec"
	"path/filepath"
)

type TranscodingCommand struct {
	Input  InputOptions
	Audio  AudioOptions
	Video  VideoOptions
	Output OutputOptions
}

type InputOptions struct {
	Format          string
	AnalyseDuration string
	ProbeSize       string
	Input           string
}

type AudioOptions struct {
	Codec      string
	SampleRate string
	AudioRate  string
}

type VideoOptions struct {
	Codec string
}

type OutputOptions struct {
	Format   string
	Time     string
	ListSize string
	Flags    string
}

// BuildCommand creates an FFmpeg command from the provided TranscodingCommand struct
func BuildCommand(tc TranscodingCommand, storagePath, streamKey string) *exec.Cmd {
	args := []string{}

	// Input options
	if tc.Input.Format != "" {
		args = append(args, "-f", tc.Input.Format)
	}
	if tc.Input.AnalyseDuration != "" {
		args = append(args, "-analyzeduration", tc.Input.AnalyseDuration)
	}
	if tc.Input.ProbeSize != "" {
		args = append(args, "-probesize", tc.Input.ProbeSize)
	}
	if tc.Input.Input != "" {
		args = append(args, "-i", tc.Input.Input)
	}

	// Audio options
	if tc.Audio.Codec != "" {
		args = append(args, "-acodec", tc.Audio.Codec)
	}
	if tc.Audio.SampleRate != "" {
		args = append(args, "-ar", tc.Audio.SampleRate)
	}
	if tc.Audio.AudioRate != "" {
		args = append(args, "-b:a", tc.Audio.AudioRate)
	}

	// Video options
	if tc.Video.Codec != "" {
		args = append(args, "-vcodec", tc.Video.Codec)
	}

	// Output options
	if tc.Output.Format != "" {
		args = append(args, "-f", tc.Output.Format)
	}
	if tc.Output.Time != "" {
		args = append(args, "-hls_time", tc.Output.Time)
	}
	if tc.Output.ListSize != "" {
		args = append(args, "-hls_list_size", tc.Output.ListSize)
	}
	if tc.Output.Flags != "" {
		args = append(args, "-hls_flags", tc.Output.Flags)
	}

	// Output path
	args = append(args, filepath.Join(storagePath, streamKey+".m3u8"))

	return exec.Command("ffmpeg", args...)
}
