package utils

type GlobalOptions struct {
    Host      string
    Port      int
    StartX    int
    StartY    int
    Loop      bool
    Threads   int

}

func NewGlobalOptions() *GlobalOptions {
	return &GlobalOptions{}
}

type ImageOptions struct {
    Path      string
    Scale     float64
    Bounce    bool
    Center    bool
    VelocityX float64
    VelocityY float64
}

func NewImageOptions() *ImageOptions {
	return &ImageOptions{}
}

type VideoOptions struct {
    Path      string
    Bounce    bool
    Center    bool
    VelocityX float64
    VelocityY float64
}

func NewVideoOptions() *VideoOptions {
	return &VideoOptions{}
}

type TextOptions struct {
    Text      string
    FontSize  float64
    FontPath  string
    Center    bool
}

func NewTextOptions() *TextOptions {
	return &TextOptions{}
}
