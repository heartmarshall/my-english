#!/usr/bin/env python3
"""
Скрипт для сбора содержимого директории и всех вложенных директорий в один txt файл.
"""

import os
from pathlib import Path
from typing import Set, List

# ==================== НАСТРОЙКИ ====================

# Директория для сканирования (путь относительно скрипта или абсолютный путь)
TARGET_DIRECTORY = "/home/alodi/playground/my-english/frontend-playground"

# Имя выходного файла
OUTPUT_FILE = "collected_content.txt"

EXCLUDED_FILENAMES: Set[str] = {
    "collected_content.txt",
    ".gitignore",
    ".DS_Store",
    "package-lock.json",
    "pnpm-lock.yaml",
    "yarn.lock",
    "bun.lockb",
    "Cargo.lock",
    "Gemfile.lock",
    "poetry.lock",
    "composer.lock",
    "go.sum",
    "collect_directory.py",
    "generated.go",
}

EXCLUDED_EXTENSIONS: Set[str] = {
    "pyc",
    "pyo",
    "__pycache__",
    "png",
    "jpg",
    "jpeg",
    "gif",
    "svg",
    "ico",
    "pdf",
    "zip",
    "tar",
    "gz",
    "node_modules",
    "venv",
    ".venv",
    "env",
    ".env",
    "dist",
    "build",
    ".next",
    ".vscode",
    ".idea",
}

# Исключаемые имена директорий (точное совпадение)
EXCLUDED_DIRECTORIES: Set[str] = {
    "docs",
    "e2e",
    ".git",
    "__pycache__",
    "node_modules",
    ".venv",
    "venv",
    "env",
    ".env",
    "dist",
    "build",
    ".next",
    ".vscode",
    ".idea",
}

# Исключаемые части путей (если путь содержит эту строку, файл будет исключен)
EXCLUDED_PATH_PARTS: Set[str] = {
    "/.git/",
    "/node_modules/",
    "/__pycache__/",
    "/.venv/",
    "/venv/",
    "/dist/",
    "/build/",
}

# Максимальный размер файла для чтения (в байтах). Файлы больше этого размера будут пропущены.
# Установите None, чтобы отключить ограничение
MAX_FILE_SIZE: int = 10 * 1024 * 1024  # 10 МБ по умолчанию

# Кодировка для чтения файлов (если None, будет использована системная кодировка)
FILE_ENCODING = "utf-8"

# Кодировка для выходного файла
OUTPUT_ENCODING = "utf-8"

# Добавлять ли разделители между файлами
ADD_SEPARATORS = True

# Символы для разделителя
SEPARATOR = "=" * 80

# ==================== КОД ====================


def should_exclude_file(file_path: Path) -> bool:
    """Проверяет, должен ли файл быть исключен."""
    # Проверка имени файла
    if file_path.name in EXCLUDED_FILENAMES:
        return True
    
    # Проверка расширения
    if file_path.suffix:
        extension = file_path.suffix[1:].lower()  # Убираем точку и приводим к нижнему регистру
        if extension in EXCLUDED_EXTENSIONS:
            return True
    
    # Проверка частей пути
    path_str = str(file_path)
    for excluded_part in EXCLUDED_PATH_PARTS:
        if excluded_part in path_str:
            return True
    
    return False


def should_exclude_directory(dir_path: Path) -> bool:
    """Проверяет, должна ли директория быть исключена."""
    if dir_path.name in EXCLUDED_DIRECTORIES:
        return True
    
    # Проверка частей пути
    path_str = str(dir_path)
    for excluded_part in EXCLUDED_PATH_PARTS:
        if excluded_part in path_str:
            return True
    
    return False


def read_file_content(file_path: Path) -> str:
    """Читает содержимое файла с обработкой ошибок."""
    try:
        # Проверка размера файла
        if MAX_FILE_SIZE is not None:
            file_size = file_path.stat().st_size
            if file_size > MAX_FILE_SIZE:
                return f"[Файл слишком большой: {file_size} байт, пропущен]"
        
        # Попытка прочитать файл
        if FILE_ENCODING:
            with open(file_path, 'r', encoding=FILE_ENCODING, errors='replace') as f:
                return f.read()
        else:
            with open(file_path, 'r', errors='replace') as f:
                return f.read()
    except UnicodeDecodeError:
        return f"[Ошибка: файл содержит бинарные данные и не может быть прочитан как текст]"
    except PermissionError:
        return f"[Ошибка: нет доступа к файлу]"
    except Exception as e:
        return f"[Ошибка при чтении файла: {str(e)}]"


def collect_directory_contents(root_dir: Path, output_file: Path) -> None:
    """Собирает содержимое всех файлов в директории и записывает в выходной файл."""
    collected_files: List[Path] = []
    skipped_files: List[Path] = []
    
    # Обход всех файлов в директории
    for current_path in root_dir.rglob('*'):
        # Пропускаем директории
        if current_path.is_dir():
            if should_exclude_directory(current_path):
                continue
            continue
        
        # Пропускаем файлы, которые нужно исключить
        if should_exclude_file(current_path):
            skipped_files.append(current_path)
            continue
        
        collected_files.append(current_path)
    
    # Сортировка файлов для предсказуемого порядка
    collected_files.sort()
    
    # Запись в выходной файл
    with open(output_file, 'w', encoding=OUTPUT_ENCODING) as out:
        # Заголовок
        out.write(f"Содержимое директории: {root_dir.absolute()}\n")
        out.write(f"Всего файлов собрано: {len(collected_files)}\n")
        out.write(f"Пропущено файлов: {len(skipped_files)}\n")
        out.write(f"\n{SEPARATOR}\n\n")
        
        # Запись содержимого каждого файла
        for file_path in collected_files:
            relative_path = file_path.relative_to(root_dir)
            
            if ADD_SEPARATORS:
                out.write(f"\n{SEPARATOR}\n")
            
            out.write(f"Файл: {relative_path}\n")
            out.write(f"Полный путь: {file_path.absolute()}\n")
            out.write(f"{SEPARATOR}\n\n")
            
            content = read_file_content(file_path)
            out.write(content)
            
            if ADD_SEPARATORS:
                out.write(f"\n\n{SEPARATOR}\n")
            else:
                out.write("\n\n")
        
        # Информация о пропущенных файлах
        if skipped_files:
            out.write(f"\n\n{SEPARATOR}\n")
            out.write("ПРОПУЩЕННЫЕ ФАЙЛЫ:\n")
            out.write(f"{SEPARATOR}\n\n")
            for skipped in sorted(skipped_files):
                out.write(f"  - {skipped.relative_to(root_dir)}\n")
    
    print(f"✓ Собрано файлов: {len(collected_files)}")
    print(f"✓ Пропущено файлов: {len(skipped_files)}")
    print(f"✓ Результат сохранен в: {output_file.absolute()}")


def main():
    """Главная функция."""
    # Определение целевой директории
    script_dir = Path(__file__).parent
    if os.path.isabs(TARGET_DIRECTORY):
        target_dir = Path(TARGET_DIRECTORY)
    else:
        target_dir = script_dir / TARGET_DIRECTORY
    
    if not target_dir.exists():
        print(f"Ошибка: директория '{target_dir}' не существует!")
        return
    
    if not target_dir.is_dir():
        print(f"Ошибка: '{target_dir}' не является директорией!")
        return
    
    # Определение выходного файла
    if os.path.isabs(OUTPUT_FILE):
        output_path = Path(OUTPUT_FILE)
    else:
        output_path = script_dir / OUTPUT_FILE
    
    print(f"Сканирование директории: {target_dir.absolute()}")
    print(f"Выходной файл: {output_path.absolute()}")
    print()
    
    collect_directory_contents(target_dir, output_path)


if __name__ == "__main__":
    main()

