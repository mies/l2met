# l2met

Convert your log stream into charts and actionable alerts in less than 1 minute
with 0 software installation.

## Log Conventions

### Counter Metric

Metrics Produced:

* app.module.function.count

```
measure="app.module.function"
```

### Sample Metric

Samples are useful for building metrics around time based functions. For instance, the elapsed duration of a function call. Or you can measure the value of an in memory resource.

Metrics Produced:

* app.module.function.min
* app.module.function.max
* app.module.function.mean
* app.module.function.median
* app.module.function.perc95
* app.module.function.perc99
* app.module.function.last
* app.module.function.count

Protocol:

```
measure="app.module.function" val=42
```

Examples:

```
measure="core.apps.get" val=1.23
```

In the previous example we have an app named **core** which has an HTTP GET endpoint named **apps** that took 1.23 seconds to execute.

## API

In addition to sending data to librato, l2met also exposes the metrics via an HTTP API. This makes it possible to build your own dashboard and visualization tooling.

There are several things to consider when interacting with l2met's API.

* Bucket
* Metric

A **bucket** is identified by: token, name, time, duration, and source. The bucket will hold the values for the identified metric. An example bucket looks like:

```
{
	"name": "some-function",
	"source": "myapp",
	"time": "1358486485",
	"duration": "minute",
	"vals": [2.18, 3.14, 32, 64, 128.0]
}
```

A **metric** is identified by: token, name, time, duration, and source. An example metric looks like:

```
{
	"name": "some-function",
	"source": "myapp",
	"time": "1358486485",
	"duration": "minute",
	"min": 2.18,
	"median": 32,
	"perc95": 64,
	"perc99": 128.0,
	"max": 128.0,
	"mean": 45.86
}
```

The primary difference between a bucket and a metric is that the bucket provides all of the raw measurement data and the metric sumarizes the data. For most cases, it will be more efficient to work with metrics. However, l2met exposes the buckets for more complicated use caeses.
