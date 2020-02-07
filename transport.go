package yandex

type Transport struct {
	// Основной цвет транспортного средства в шестнадцатеричном формате.
	Color string `json:"color"`
	// Код подтипа транспорта для типа, указанного в элементе transport_type. Подтип может совпадать с типом (например, для обычной электрички указывается тип suburban и подтип suburban).
	//
	// Другие возможные значения:
	//
	// helicopter — вертолет (для типа plane);
	// rex — экспресс РЭКС (для типа suburban);
	// sputnik — «Спутник» (для типа suburban);
	// skiarrow — «Лыжная стрела» (для типа suburban);
	// shezh — «Снежинка» (для типа suburban);
	// skirus — «Лыжня России» (для типа suburban);
	// city — городская электричка (для типа suburban);
	// kalina — «Калина красная» (для типа suburban);
	// vostok — «Восток» (для типа suburban);
	// prostoryaltaya — «Просторы Алтая» (для типа suburban);
	// 14vag — состав из 14 вагонов (для типа suburban);
	// last — «Ласточка» (для типа suburban);
	// exprdal — экспресс с билетами на конкретные места (для типа suburban);
	// volzhex — «Волжский экспресс» (для типа suburban);
	// stdplus — электрички типа «стандарт плюс» (для типа suburban);
	// express — экспресс (для типа suburban);
	// skor — ускоренный поезд (для типа suburban);
	// fiztekh — Физтех.Электричка (для типа suburban);
	// vag6 — состав из 6 вагонов (для типа suburban);
	// river— речной транспорт (для типа water);
	// sea — морской транспорт (для типа water).
	Code string `json:"code"`
	// Описание подтипа транспорта на естественном языке.
	Title string `json:"title"`
}
