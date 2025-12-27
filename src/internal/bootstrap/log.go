package bootstrap

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/OpenListTeam/OpenList/v4/cmd/flags"
	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

type TagLogHooker struct {
	Writer    io.Writer
	Tags      []string
	LevelsArr []logrus.Level
}

func (h *TagLogHooker) Levels() []logrus.Level {
	return h.LevelsArr
}

func (h *TagLogHooker) Fire(entry *logrus.Entry) error {

	if h.Writer == nil {
		return nil
	}

	line, _ := entry.String()

	for _, k := range h.Tags {
		if strings.Contains(line, fmt.Sprintf("[%v]", k)) {
			_, err := h.Writer.Write([]byte(line))
			return err
		}
	}
	return nil
}

func createTagLogHooker(tag string) *TagLogHooker {
	hooker := new(TagLogHooker)

	logConfig := conf.Conf.Log
	if !logConfig.Enable {
		return nil
	}

	tagLogName := path.Join(filepath.Dir(logConfig.Name), fmt.Sprintf("%v_%v", strings.ToLower(tag), filepath.Base(logConfig.Name)))

	hooker.Writer = &lumberjack.Logger{
		Filename:   tagLogName,
		MaxSize:    logConfig.MaxSize, // megabytes
		MaxBackups: logConfig.MaxBackups,
		MaxAge:     logConfig.MaxAge,   //days
		Compress:   logConfig.Compress, // disabled by default
	}

	hooker.Tags = append(hooker.Tags, tag)
	hooker.LevelsArr = logrus.AllLevels
	return hooker
}

func init() {
	formatter := logrus.TextFormatter{
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		TimestampFormat:           "2006-01-02 15:04:05",
		FullTimestamp:             true,
	}
	logrus.SetFormatter(&formatter)
	utils.Log.SetFormatter(&formatter)
	// logrus.SetLevel(logrus.DebugLevel)
}

func setLog(l *logrus.Logger) {
	if flags.Debug || flags.Dev {
		l.SetLevel(logrus.DebugLevel)
		l.SetReportCaller(true)
	} else {
		l.SetLevel(logrus.InfoLevel)
		l.SetReportCaller(false)
	}
}

func Log() {
	setLog(logrus.StandardLogger())
	setLog(utils.Log)
	logConfig := conf.Conf.Log
	if logConfig.Enable {
		var w io.Writer = &lumberjack.Logger{
			Filename:   logConfig.Name,
			MaxSize:    logConfig.MaxSize, // megabytes
			MaxBackups: logConfig.MaxBackups,
			MaxAge:     logConfig.MaxAge,   //days
			Compress:   logConfig.Compress, // disabled by default
		}
		if flags.Debug || flags.Dev || flags.LogStd {
			w = io.MultiWriter(os.Stdout, w)
		}
		logrus.SetOutput(w)

		accessLogHooker := createTagLogHooker("ACCESS")
		if accessLogHooker != nil {
			logrus.AddHook(accessLogHooker)
		}

		logrus.Info("[ACCESS]access log init")

	}

	log.SetOutput(logrus.StandardLogger().Out)
	utils.Log.Infof("init logrus...")
	utils.Log = logrus.StandardLogger()

}
