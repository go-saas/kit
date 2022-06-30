package pulsar

import (
	"github.com/apache/pulsar-client-go/pulsar/log"
	klog "github.com/go-kratos/kratos/v2/log"
)

type logger struct {
	l klog.Logger
}

var _ log.Logger = (*logger)(nil)

func (l logger) SubLogger(fields log.Fields) log.Logger {
	var kv []interface{}
	for k, v := range fields {
		kv = append(kv, k, v)
	}
	return logger{l: klog.With(l.l, kv...)}
}

func (l logger) WithFields(fields log.Fields) log.Entry {
	var kv []interface{}
	for k, v := range fields {
		kv = append(kv, k, v)
	}
	return logEntry{l: klog.With(l.l, kv...)}
}

func (l logger) WithField(name string, value interface{}) log.Entry {
	return logEntry{l: klog.With(l.l, name, value)}
}

func (l logger) WithError(err error) log.Entry {
	return logEntry{l: klog.With(l.l, "error", err)}
}

func (l logger) Debug(args ...interface{}) {
	klog.NewHelper(l.l).Debug(args)
}

func (l logger) Info(args ...interface{}) {
	klog.NewHelper(l.l).Info(args)
}

func (l logger) Warn(args ...interface{}) {
	klog.NewHelper(l.l).Warn(args)
}

func (l logger) Error(args ...interface{}) {
	klog.NewHelper(l.l).Error(args)
}

func (l logger) Debugf(format string, args ...interface{}) {
	klog.NewHelper(l.l).Debugf(format, args)
}

func (l logger) Infof(format string, args ...interface{}) {
	klog.NewHelper(l.l).Infof(format, args)
}

func (l logger) Warnf(format string, args ...interface{}) {
	klog.NewHelper(l.l).Warnf(format, args)
}

func (l logger) Errorf(format string, args ...interface{}) {
	klog.NewHelper(l.l).Errorf(format, args)
}

type logEntry struct {
	l klog.Logger
}

func (l logEntry) WithFields(fields log.Fields) log.Entry {
	var kv []interface{}
	for k, v := range fields {
		kv = append(kv, k, v)
	}
	return logger{l: klog.With(l.l, kv...)}
}

func (l logEntry) WithField(name string, value interface{}) log.Entry {
	return logEntry{l: klog.With(l.l, name, value)}
}

func (l logEntry) Debug(args ...interface{}) {
	klog.NewHelper(l.l).Debug(args)
}

func (l logEntry) Info(args ...interface{}) {
	klog.NewHelper(l.l).Info(args)
}

func (l logEntry) Warn(args ...interface{}) {
	klog.NewHelper(l.l).Warn(args)
}

func (l logEntry) Error(args ...interface{}) {
	klog.NewHelper(l.l).Error(args)
}

func (l logEntry) Debugf(format string, args ...interface{}) {
	klog.NewHelper(l.l).Debugf(format, args)
}

func (l logEntry) Infof(format string, args ...interface{}) {
	klog.NewHelper(l.l).Infof(format, args)
}

func (l logEntry) Warnf(format string, args ...interface{}) {
	klog.NewHelper(l.l).Warnf(format, args)
}

func (l logEntry) Errorf(format string, args ...interface{}) {
	klog.NewHelper(l.l).Errorf(format, args)
}

var _ log.Entry = (*logEntry)(nil)
