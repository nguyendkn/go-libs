module github.com/nguyendkn/go-libs/rtsp

go 1.24

require (
	github.com/nguyendkn/go-libs/ffmpeg v0.0.0
	github.com/nguyendkn/go-libs/hls v0.0.0
)

replace (
	github.com/nguyendkn/go-libs/ffmpeg => ../ffmpeg
	github.com/nguyendkn/go-libs/hls => ../hls
)
