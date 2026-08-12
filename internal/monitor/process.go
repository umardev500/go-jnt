package monitor

import (
	"github.com/rs/zerolog"
	"github.com/xuri/excelize/v2"
)

func Process(
	input string,
	output string,
	sheet string,
	lookupFile string,
	log zerolog.Logger,
) error {

	f, err := excelize.OpenFile(input)

	if err != nil {
		return err
	}

	defer f.Close()

	err = RemoveUnusedColumns(
		f,
		sheet,
		[]string{
			"M",
			"L",
			"J",
			"I",
			"G",
			"E",
			"D",
		},
		log,
	)

	if err != nil {
		return err
	}

	err = InsertStyledColumnsAfterC(
		f,
		sheet,
		log,
	)

	if err != nil {
		return err
	}

	err = UpdateTujuanFromRetur(
		f,
		sheet,
		log,
	)

	if err != nil {
		return err
	}

	err = FillKlasifikasi(
		f,
		sheet,
		lookupFile,
		"BTN KE TUJUAN",
		log,
	)

	if err != nil {
		return err
	}

	err = FillProvinsi(
		f,
		sheet,
		lookupFile,
		"PROVINSI ESTIMASI BARANG",
		log,
	)

	// Update JKT_GATEWAY based on Jenis Layanan
	if err := UpdateKlasifikasiByJenisLayanan(
		f,
		sheet,
		log,
	); err != nil {
		log.Fatal().
			Err(err).
			Msg("failed updating klasifikasi by jenis layanan")
	}

	if err := UpdateIncomingByLokasiSebelumnya(
		f,
		sheet,
		log,
	); err != nil {
		log.Fatal().
			Err(err).
			Msg("failed updating incoming")
	}

	if err != nil {
		return err
	}

	if err := UpdateBelumPickup(
		f,
		sheet,
		log,
	); err != nil {
		log.Fatal().
			Err(err).
			Msg("failed updating belum pickup")
	}

	// if err := UpdateJKTIfKlasifikasiNA(
	// 	f,
	// 	sheet,
	// 	log,
	// ); err != nil {
	// 	log.Fatal().
	// 		Err(err).
	// 		Msg("failed updating JKT for #N/A klasifikasi")
	// }

	return f.SaveAs(output)
}
