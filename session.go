package gompv

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

var SessionAlreadyStartedError = fmt.Errorf("session already started")

type Source struct {
	ID   string
	Path string
}

type Session struct {
	client     *MPVClient
	command    *Command
	CusomFlags map[string]string // Custom flags for the session
	cmd        *exec.Cmd
	videos     []*Source
	audios     []*Source
	filters    []*Filter
}

// NewSession creates a new MPV session with the default command.
// It initializes the MPV client and command, but does not start the process.
// You can add flags to the command before starting the session.
func NewSession(videoPath string, audioPaths ...string) *Session {
	s := &Session{
		client:     NewMPVClient(),
		CusomFlags: make(map[string]string),
	}
	s.AddVideoSource(videoPath)
	for _, audioPath := range audioPaths {
		s.AddAudioSource(audioPath)
	}
	return s
}

func NewSessionWithSocketPath(socketPath string) *Session {
	return &Session{
		client:     NewMPVClientWithSocketPath(socketPath),
		CusomFlags: make(map[string]string),
	}
}

func (s *Session) AddVideoSource(path string) error {
	if s.cmd != nil {
		return SessionAlreadyStartedError
	}
	id := fmt.Sprintf("vid%d", len(s.videos)+1) // Unique ID for each video
	s.videos = append(s.videos, &Source{ID: id, Path: path})
	return nil
}

func (s *Session) AddAudioSource(path string) error {
	if s.cmd != nil {
		return SessionAlreadyStartedError
	}
	no := len(s.audios)
	id := fmt.Sprintf("aid%d", no+1) // Unique ID for each audio
	s.audios = append(s.audios, &Source{ID: id, Path: path})
	return nil
}

func (s *Session) AddCustomFilter(filter *Filter) error {
	if s.cmd != nil {
		return SessionAlreadyStartedError
	}
	s.filters = append(s.filters, filter)
	return nil
}

func (s *Session) AddGlobalAudioFilter(operator string) error {
	if s.cmd != nil {
		return SessionAlreadyStartedError
	}
	if len(s.audios) == 0 {
		return fmt.Errorf("no audio sources available to add global audio filter")
	}
	filter := NewFilter(operator)
	// Set the target of the filter to "ao" for global audio processing
	filter.SetTarget("ao")
	for _, audio := range s.audios {
		filter.In = append(filter.In, audio.ID) // Add all audio sources to the AddGlobalAudioFilter
	}
	s.filters = append(s.filters, filter)
	return nil
}

func (s *Session) Command() *Command {
	if s.command == nil {
		s.command = NewCommand()
	}
	if s.command.Flags["input-ipc-server"] == "" {
		s.command.AddFlag("input-ipc-server", s.client.SocketPath())
	}
	for flag, value := range s.CusomFlags {
		s.command.AddFlag(flag, value)
	}

	audioIds := make([]string, len(s.audios))
	// Add audio sources to the command
	for i, audio := range s.audios {
		audioIds[i] = audio.ID
		if i == 0 {
			s.command.AddFlag(
				"audio-files",
				audio.Path,
			) // Use the first audio file as the main audio
		} else {
			s.command.AddFlag("audio-files-append", audio.Path) // Add additional audio files
		}
	}

	videoIds := make([]string, len(s.videos))
	// Add video sources to the command
	for i, video := range s.videos {
		videoIds[i] = video.ID
		s.command.AddArg(video.Path)
	}
	// Add filters to the command
	for _, filter := range s.filters {
		s.command.AddFlag("lavfi-complex", fmt.Sprintf("'%s'", filter.String()))
	}
	return s.command
}

func (s *Session) Start() error {
	s.Command() // Ensure the command is initialized
	fmt.Println("Starting MPV with command:", s.command.String())
	s.cmd = exec.Command("bash", "-c", s.command.String())
	s.cmd.Stdout = os.Stdout
	s.cmd.Stderr = os.Stderr
	err := s.cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start MPV process: %w", err)
	}
	time.Sleep(1 * time.Second) // Give some time for the command to be prepared
	err = s.client.Open()
	if err != nil {
		return fmt.Errorf("failed to open MPV client: %w", err)
	}
	go func() {
		s.client.WaitUntilClosed()
		s.cmd.Process.Kill() // Ensure the process is killed when the client closes
	}()
	go func() {
		s.cmd.Wait()
		s.client.Close() // Close the client when the process exits
	}()
	return nil
}

func (s *Session) Stop() error {
	if s.cmd != nil {
		clientErr := s.client.Close()
		processErr := s.cmd.Process.Kill()
		if clientErr != nil && processErr != nil {
			return fmt.Errorf("failed to stop session: %v, %v", clientErr, processErr)
		}
		if clientErr != nil {
			return fmt.Errorf("failed to close MPV client: %w", clientErr)
		}
		if processErr != nil {
			return fmt.Errorf("failed to kill process: %w", processErr)
		}
		s.cmd = nil
		s.client = nil
		return nil
	}
	return nil
}

func (s *Session) MPVCommand() *exec.Cmd {
	if s.cmd == nil {
		s.Command() // Ensure the command is initialized
		s.cmd = exec.Command("bash", "-c", s.command.String())
		s.cmd.Stdout = os.Stdout
		s.cmd.Stderr = os.Stderr
	}
	return s.cmd
}

func (s *Session) SetInputIPCSocket(socketPath string) error {
	if s.cmd != nil {
		return SessionAlreadyStartedError
	}
	s.client = NewMPVClientWithSocketPath(socketPath)
	return nil
}

func (s *Session) Client() *MPVClient {
	return s.client
}

func (s *Session) AddFlag(flag, value string) {
	s.CusomFlags[flag] = value
}

func (s *Session) SetOSC(value bool) {
	if value {
		s.CusomFlags["osc"] = "yes"
	} else {
		s.CusomFlags["osc"] = "no"
	}
}

func (s *Session) SetNoInputDefaultBindings(value bool) {
	if value {
		s.CusomFlags["no-input-default-bindings"] = ""
	}
}
