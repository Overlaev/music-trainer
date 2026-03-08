// 1. Настройки игры
const NOTES =[
    "C4", "C#4", "D4", "D#4", "E4", "F4", "F#4", "G4", "G#4", "A4", "A#4", "B4",
    "C5", "C#5", "D5", "D#5", "E5", "F5", "F#5", "G5", "G#5", "A5", "A#5", "B5"
];
const GAME_DURATION = 60; 
const MAX_POINTS_PER_NOTE = 1000;
const MAX_TIME_FOR_MAX_POINTS = 3000; 
const PENALTY_POINTS = 100;

let score = 0;
let timeLeft = GAME_DURATION;
let currentNote = "";
let noteAppearedAt = 0;
let gameInterval;
let isPlaying = false;

// Элементы DOM
const scoreEl = document.getElementById("score");
const timerEl = document.getElementById("timer");
const targetNoteEl = document.getElementById("target-note");
const pianoEl = document.getElementById("piano");
const startBtn = document.getElementById("start-btn");

// ИНИЦИАЛИЗАЦИЯ ЗВУКА
const pianoSampler = new Tone.Sampler({
    urls: {
        "C4": "C4.mp3",
        "D#4": "Ds4.mp3",
        "F#4": "Fs4.mp3",
        "A4": "A4.mp3",
    },
    release: 1,
    baseUrl: "https://tonejs.github.io/audio/salamander/",
}).toDestination();

// 2. Инициализация пианино (ИСПРАВЛЕННАЯ: теперь задаем класс .white)
function initPiano() {
    pianoEl.innerHTML = "";
    NOTES.forEach(note => {
        const btn = document.createElement("div"); // Используем div
        btn.className = "key";
        
        // Разделяем на черные и белые клавиши
        if (note.includes("#")) {
            btn.classList.add("black");
            btn.innerText = note.replace("#", "♯"); 
        } else {
            btn.classList.add("white");
            btn.innerText = note; 
        }
        
        btn.onclick = () => handleKeyPress(note, btn);
        pianoEl.appendChild(btn);
    });
}

// 3. ОТРИСОВКА НОТ (VexFlow)
function drawNoteOnStaff(noteString) {
    targetNoteEl.innerHTML = ""; 
    
    const VF = Vex.Flow;
    const renderer = new VF.Renderer(targetNoteEl, VF.Renderer.Backends.SVG);
    renderer.resize(150, 150);
    const context = renderer.getContext();

    const stave = new VF.Stave(10, 40, 130);
    stave.addClef("treble").setContext(context).draw();

    const letter = noteString.charAt(0).toLowerCase();
    const accidental = noteString.includes("#") ? "#" : "";
    const octave = noteString.charAt(noteString.length - 1);
    const vexKey = `${letter}${accidental}/${octave}`;

    const staveNote = new VF.StaveNote({ clef: "treble", keys:[vexKey], duration: "w" });
    
    if (accidental) {
        staveNote.addModifier(new VF.Accidental(accidental));
    }

    const voice = new VF.Voice({ num_beats: 4, beat_value: 4 });
    voice.addTickables([staveNote]);
    new VF.Formatter().joinVoices([voice]).format([voice], 80);
    voice.draw(context, stave);
}

// 4. Старт игры
startBtn.onclick = async () => {
    await Tone.start();
    
    score = 0;
    timeLeft = GAME_DURATION;
    isPlaying = true;
    startBtn.style.display = "none";
    scoreEl.innerText = `Очки: ${score}`;
    
    nextNote();
    
    gameInterval = setInterval(() => {
        timeLeft--;
        timerEl.innerText = `Время: ${timeLeft}`;
        if (timeLeft <= 0) endGame();
    }, 1000);
};

// 5. Следующая нота
function nextNote() {
    let newNote;
    do {
        newNote = NOTES[Math.floor(Math.random() * NOTES.length)];
    } while (newNote === currentNote);
    
    currentNote = newNote;
    drawNoteOnStaff(currentNote); 
    noteAppearedAt = Date.now();
}

// 6. Обработка нажатий
function handleKeyPress(clickedNote, btnElement) {
    if (!isPlaying) return;

    pianoSampler.triggerAttackRelease(clickedNote, "8n");

    if (clickedNote === currentNote) {
        // Правильно
        const timeTaken = Date.now() - noteAppearedAt;
        let earnedPoints = Math.max(100, Math.round(MAX_POINTS_PER_NOTE * (1 - timeTaken / MAX_TIME_FOR_MAX_POINTS)));
        if (timeTaken > MAX_TIME_FOR_MAX_POINTS) earnedPoints = 100;

        score += earnedPoints;
        scoreEl.innerText = `Очки: ${score}`;
        
        btnElement.classList.add("correct");
        setTimeout(() => {
            btnElement.classList.remove("correct");
            nextNote();
        }, 200);

    } else {
        // Неправильно
        score = Math.max(0, score - PENALTY_POINTS); 
        scoreEl.innerText = `Очки: ${score}`;

        btnElement.classList.add("wrong");
        setTimeout(() => {
            btnElement.classList.remove("wrong");
        }, 300);
    }
}

// 7. Конец игры (Связь с Golang Backend)
async function endGame() {
    isPlaying = false;
    clearInterval(gameInterval);
    targetNoteEl.innerHTML = "<div style='font-size:30px; margin-top:50px'>Время вышло!</div>";
    startBtn.style.display = "inline-block";
    startBtn.innerText = "Играть снова";

    const username = prompt(`Игра окончена!\nТвой счет: ${score} очков.\nВведите имя для сохранения результата:`, "Игрок");

    if (username !== null) {
        try {
            await fetch('http://localhost:8080/api/scores', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username: username, score: score })
            });

            await loadLeaderboard();
            
        } catch (error) {
            console.error("Ошибка при отправке на сервер:", error);
            alert("Не удалось подключиться к серверу (убедитесь, что Go-сервер запущен)");
        }
    }
}

// 8. Загрузка Таблицы Лидеров
async function loadLeaderboard() {
    try {
        const response = await fetch('http://localhost:8080/api/leaderboard');
        const topScores = await response.json();

        let html = "<h3>🏆 Таблица лидеров</h3><ol style='text-align:left; display:inline-block;'>";
        topScores.forEach(entry => {
            html += `<li style="margin-bottom: 5px;">
                        <b>${entry.username}</b>: ${entry.score} очков
                     </li>`;
        });
        html += "</ol>";

        targetNoteEl.innerHTML = html;

    } catch (error) {
        console.error("Ошибка при получении таблицы:", error);
    }
}

// ЗАПУСК Отрисовки пианино при загрузке страницы
initPiano();