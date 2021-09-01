package main

import (
	"github.com/samuelventura/laurelview/pkg/lvsdk"
)

var M = lvsdk.M
var Mn = lvsdk.Mn
var Mna = lvsdk.Mna
var Mns = lvsdk.Mns
var Mnsa = lvsdk.Mnsa
var NewId = lvsdk.NewId
var NewCount = lvsdk.NewCount
var NewEntry = lvsdk.NewEntry
var NewCleaner = lvsdk.NewCleaner
var NewRuntime = lvsdk.NewRuntime
var NewSocketDial = lvsdk.NewSocketDial
var NewSocketConn = lvsdk.NewSocketConn
var NewTestOutput = lvsdk.NewTestOutput
var DefaultRuntime = lvsdk.DefaultRuntime

var NopAction = lvsdk.NopAction
var NopDispatch = lvsdk.NopDispatch
var NopOutput = lvsdk.NopOutput
var NopFactory = lvsdk.NopFactory
var NopLogger = lvsdk.NopLogger
var NopHandler = lvsdk.NopHandler

var AsyncDispatch = lvsdk.AsyncDispatch
var ClearDispatch = lvsdk.ClearDispatch
var MapDispatch = lvsdk.MapDispatch
var DisposeArgs = lvsdk.DisposeArgs

var AssertTrue = lvsdk.AssertTrue
var PanicIfError = lvsdk.PanicIfError
var TraceRecover = lvsdk.TraceRecover
var TraceIfError = lvsdk.TraceIfError
var ErrorString = lvsdk.ErrorString
var PanicLN = lvsdk.PanicLN
var PanicF = lvsdk.PanicF

var CloseLog = lvsdk.CloseLog
var DefaultLog = lvsdk.DefaultLog
var DefaultOutput = lvsdk.DefaultOutput
var LevelOutput = lvsdk.LevelOutput
var PrefixLogger = lvsdk.PrefixLogger
var SimpleLogger = lvsdk.SimpleLogger
var FlatPrintln = lvsdk.FlatPrintln
var GoLogLogger = lvsdk.GoLogLogger

var WaitClose = lvsdk.WaitClose
var WaitChannel = lvsdk.WaitChannel
var SendChannel = lvsdk.SendChannel

var Future = lvsdk.Future
var Millis = lvsdk.Millis
var Readable = lvsdk.Readable

var EnvironDefault = lvsdk.EnvironDefault
var EnvironFromFile = lvsdk.EnvironFromFile
var RelativeSibling = lvsdk.RelativeSibling
var RelativeExtension = lvsdk.RelativeExtension
var ChangeExtension = lvsdk.ChangeExtension
var ExecutablePath = lvsdk.ExecutablePath
var WaitExitSignal = lvsdk.WaitExitSignal

var ReadMutation = lvsdk.ReadMutation
var WriteMutation = lvsdk.WriteMutation
var EncodeMutation = lvsdk.EncodeMutation
var DecodeMutation = lvsdk.DecodeMutation
var ParseUint = lvsdk.ParseUint
var ParseString = lvsdk.ParseString
var MaybeUint = lvsdk.MaybeUint
var MaybeString = lvsdk.MaybeString
var CastUint = lvsdk.CastUint

type Handler = lvsdk.Handler
type Entry = lvsdk.Entry
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
