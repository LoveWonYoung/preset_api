"""把仓库 python/ 目录加入 sys.path，供示例脚本导入绑定。"""

from __future__ import annotations

import sys
from pathlib import Path

_PYTHON_DIR = Path(__file__).resolve().parents[2] / "python"
_path = str(_PYTHON_DIR)
if _path not in sys.path:
    sys.path.insert(0, _path)
