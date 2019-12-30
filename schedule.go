package yandex

import (
	"time"
)

type Schedule struct {
	ExceptDays string    `json:"except_days"` // Дни, в которые нитка не курсирует (даже если они входят в множество, описанное элементом days). Format "6, 7, 8, 9, 13, 14 февраля"
	Arrival    time.Time `json:"arrival"`     // Время прибытия
	Thread     Thread    `json:"thread"`      // Информация о нитке
	IsFuzzy    bool      `json:"is_fuzzy"`    // Признак неточности времени отправления и времени прибытия. Возможные значения: true — время прибытия и время отправления указаны неточно; false — время прибытия и время отправления указан точно.
	Days       string    `json:"days"`        // Дни курсирования нитки
	Stops      string    `json:"stops"`       // Станции следования рейса, на которых совершается остановка. Описывается в свободной форме. Например, значение везде значит, что остановка совершается на всех станциях следования. Пустая строка значит, что нитка нигде не останавливается между начальной и конечной станциями.
	Departure  time.Time `json:"departure"`   // Время отправления
	Terminal   string    `json:"terminal"`    // Терминал аэропорта (например, «D»). Принимает значение null, если информации о терминале нет.
	Platform   string    `json:"platform"`    // Платформа или путь, с которого отправляется рейс (например, «3 путь»). Пустая строка значит, что информации о платформе или пути нет.
}
