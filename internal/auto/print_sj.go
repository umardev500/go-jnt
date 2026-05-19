package auto

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func setPrinter(name string) {
	exec.Command(
		"rundll32",
		"printui.dll,PrintUIEntry",
		"/y",
		"/n",
		name,
	).Run()
}

func PrintSJ(path string) {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// =========================
	// 0. Save current printer (Zebra)
	// =========================
	originalPrinter := "ZDesigner ZD220-203dpi ZPL"

	// =========================
	// 1. Switch to HP printer
	// =========================
	setPrinter("HP Laser MFP 131 133 135-139")

	// small delay (important for Windows)
	// time.Sleep(2 * time.Second)

	// =========================
	// 2. Start Excel
	// =========================
	// ole.CoInitialize(0)
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED)
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("Excel.Application")
	if err != nil {
		panic(err)
	}

	excel, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		panic(err)
	}
	defer excel.Release()

	oleutil.PutProperty(excel, "Visible", true)

	// =========================
	// 3. Open workbook
	// =========================
	workbooks := oleutil.MustGetProperty(excel, "Workbooks").ToIDispatch()
	// fileToPrint := config.GetSuratJalanPrintFilePath()
	wb := oleutil.MustCallMethod(workbooks, "Open", path).ToIDispatch()

	// =========================
	// 4. Run macro
	// =========================
	macro := "'" + path + "'!Module2.PrintWithArea"

	_, err = oleutil.CallMethod(excel, "Run", macro)
	if err != nil {
		panic("macro failed: " + err.Error())
	}
	// oleutil.MustCallMethod(excel, "Run", "Module2.PrintWithArea")

	// =========================
	// 5. Save & close
	// =========================
	oleutil.MustCallMethod(wb, "Save")
	oleutil.MustCallMethod(wb, "Close", false)

	oleutil.MustCallMethod(excel, "Quit")

	// =========================
	// 6. Restore Zebra printer
	// =========================
	setPrinter(originalPrinter)

	fmt.Println("Macro executed")
}
