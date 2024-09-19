package log

import (
	"MagicWand/library/log/internal/core"
	"math"
	"runtime"
	"strconv"
	"time"
)

func funcName(skip int) (name string) {
	if _, file, lineNo, ok := runtime.Caller(skip); ok {
		return file + ":" + strconv.Itoa(lineNo)
	}
	return "unknown:0"
}

// toMap convert D slice to map[string]interface{} for legacy file and stdout.
func toMap(args ...D) map[string]interface{} {
	d := make(map[string]interface{}, 10+len(args))
	for _, arg := range args {
		switch arg.Type {
		case core.UintType, core.Uint64Type, core.IntTpye, core.Int64Type:
			d[arg.Key] = arg.Int64Val
		case core.StringType:
			d[arg.Key] = arg.StringVal
		case core.Float32Type:
			d[arg.Key] = math.Float32frombits(uint32(arg.Int64Val))
		case core.Float64Type:
			d[arg.Key] = math.Float64frombits(uint64(arg.Int64Val))
		case core.DurationType:
			d[arg.Key] = time.Duration(arg.Int64Val)
		default:
			d[arg.Key] = arg.Value
		}
	}
	return d
}

//func addExtraField(ctx context.Context, fields map[string]interface{}) {
//	if t, ok := trace.FromContext(ctx); ok {
//		traceFlags := "00"
//		if t.IsSampled() {
//			traceFlags = "01"
//		}
//		fields[_tid] = t.TraceID()
//		fields[_span] = t.SpanID()
//		fields[_traceFlags] = traceFlags
//	}
//	if caller := metadata.String(ctx, metadata.Caller); caller != "" {
//		fields[_caller] = caller
//	}
//	if color := metadata.String(ctx, metadata.Color); color != "" {
//		fields[_color] = color
//	}
//	if env.Color != "" {
//		fields[_envColor] = env.Color
//	}
//	if cluster := metadata.String(ctx, metadata.Cluster); cluster != "" {
//		fields[_cluster] = cluster
//	}
//	fields[_deplyEnv] = env.DeployEnv
//	fields[_zone] = env.Zone
//	c := c()
//	fields[_appID] = c.Family
//	fields[_instanceID] = c.Host
//	if mirror := metadata.String(ctx, metadata.Mirror); mirror != "" {
//		fields[_mirror] = mirror
//	}
//	tenant, ok := tenant.FromContext(ctx)
//	if ok {
//		fields[_tenantKey] = tenant.TenantKey
//	}
//}
