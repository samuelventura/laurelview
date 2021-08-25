package main

import (
	"github.com/samuelventura/laurelview/pkg/lvnrt"
	"github.com/samuelventura/laurelview/pkg/lvsdk"
)

var NewBus = lvnrt.NewBus
var NewHub = lvnrt.NewHub
var NewState = lvnrt.NewState
var NewEntry = lvsdk.NewEntry

var M = lvsdk.M
var Mn = lvsdk.Mn
var Mna = lvsdk.Mna
var Mns = lvsdk.Mns
var NewId = lvsdk.NewId
var NewDpm = lvnrt.NewDpm
var NewCount = lvsdk.NewCount
var NewSocket = lvsdk.NewSocket
var NewCleaner = lvsdk.NewCleaner
var NewRuntime = lvsdk.NewRuntime
var NewSocketConn = lvsdk.NewSocketConn
var NewTestOutput = lvsdk.NewTestOutput
var DefaultRuntime = lvsdk.DefaultRuntime

var NopAction = lvsdk.NopAction
var NopDispatch = lvsdk.NopDispatch
var NopOutput = lvsdk.NopOutput
var NopFactory = lvsdk.NopFactory

var AsyncDispatch = lvsdk.AsyncDispatch
var ClearDispatch = lvsdk.ClearDispatch
var MapDispatch = lvsdk.MapDispatch
var DisposeArgs = lvsdk.DisposeArgs

var AssertTrue = lvsdk.AssertTrue
var PanicIfError = lvsdk.PanicIfError
var TraceRecover = lvsdk.TraceRecover
var TraceIfError = lvsdk.TraceIfError
var PanicLN = lvsdk.PanicLN
var PanicF = lvsdk.PanicF

var CloseLog = lvsdk.CloseLog
var DefaultLog = lvsdk.DefaultLog
var DefaultOutput = lvsdk.DefaultOutput
var LevelOutput = lvsdk.LevelOutput
var PrefixLogger = lvsdk.PrefixLogger

var WaitClose = lvsdk.WaitClose

var Future = lvsdk.Future
var Millis = lvsdk.Millis

type Log = lvsdk.Log
type Output = lvsdk.Output
type Map = lvsdk.Map
type Queue = lvsdk.Queue
type Channel = lvsdk.Channel
type Any = lvsdk.Any
type Action = lvsdk.Action
type Dispatch = lvsdk.Dispatch
type Factory = lvsdk.Factory
type Runtime = lvsdk.Runtime
type Mutation = lvsdk.Mutation
type Logger = lvsdk.Logger
type Socket = lvsdk.Socket
type Cleaner = lvsdk.Cleaner
type Count = lvsdk.Count
type Id = lvsdk.Id
type TestOutput = lvsdk.TestOutput
