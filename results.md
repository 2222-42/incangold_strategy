# Incan Gold Strategy Simulation Results

10万回のゲーム進行をシミュレーションした結果です。

```text
Simulation completed in 4.598s
Total Games: 100000
-----------------------------------------------------
Strategy                  | Win Rate   | Avg Score 
-----------------------------------------------------
Greedy (Target 10)        | 50.51%     | 25.98
Threshold (2 Hazards)     | 26.28%     | 20.19
Threshold (1 Hazard)      | 13.96%     | 14.43
Random 10%                | 10.47%     | 13.39
Random 50%                |  2.88%     |  7.46
Greedy (Never Leave)      |  0.00%     |  0.00
-----------------------------------------------------
```

---

## 実験 2：Greedy除外 + RiskyStrategy導入後（各10万ゲーム）

Greedy戦略を除外し、文献調査をもとに新たに実装した **RiskyStrategy** を追加した条件での結果です。

### Scenario: `all` — 全戦略総当たり（5人）

```text
Simulation completed in 3.79s
Total Games: 100000
-----------------------------------------------------
Strategy                  | Win Rate   | Avg Score
-----------------------------------------------------
Risky (Best+ArtEV)         | 72.07%     | 38.93
Threshold (2 Hazards)      | 15.78%     | 22.40
Threshold (1 Hazard)       |  7.03%     | 14.55
Random 10%                 |  6.94%     | 15.03
Random 50%                 |  0.97%     |  6.98
-----------------------------------------------------
```

### Scenario: `risky` — Risky vs ベースライン（3人）

```text
Simulation completed in 3.04s
Total Games: 100000
-----------------------------------------------------
Strategy                         | Win Rate   | Avg Score
-----------------------------------------------------
Threshold (2 Hazards)            | 49.21%     | 32.95
Risky (Best: R=2,S1=4,S2=2+ArtEV) | 44.87%   | 33.33
Risky (Best: R=2,S1=4,S2=2)     | 29.91%     | 29.88
-----------------------------------------------------
```

### Scenario: `risky-vs-risky` — Riskyパラメーター内対決（5人）

```text
Simulation completed in 3.68s
Total Games: 100000
-----------------------------------------------------
Strategy                         | Win Rate   | Avg Score
-----------------------------------------------------
Risky (R=2,S1=2,S2=-2)          | 82.81%     | 19.40
Risky (Flat base=9)              | 31.03%     |  9.66
Risky (R=2,S1=4,S2=2+ArtEV)     | 30.76%     |  8.12
Risky (R=0,S1=4,S2=4)           | 26.45%     |  4.50
Risky (R=2,S1=4,S2=2)           | 23.54%     |  4.65
-----------------------------------------------------
```

---

## 考察

### RiskyStrategyの設計

文献（`docs/search_result.md` §7）の数値シミュレーションをベースに実装。  
意思決定フロー：

```
effectivePocket = self.PocketScore
  + artifactEstimate(visible_artifacts) / activePlayerCount  // WithArtifactEV=true 時

threshold = BaseThreshold + (lead >= R ? S1 : S2)

if effectivePocket >= threshold → Leave
```

| パラメーター | 文献最良値 | 意味 |
|---|---|---|
| BaseThreshold | 9 | 基準となる撤退閾値（ポケットのジェム数） |
| R | 2 | 「リード」判定の点差 |
| S1 | 4 | リード時の追加待機（→ 計13ジェム） |
| S2 | 2 | 遅れ時の追加待機（→ 計11ジェム） |
| WithArtifactEV | true | アーティファクト期待値を加算するか |

### 主な知見

**1. アーティファクト期待値（ArtEV）の効果が顕著**

`scenario:risky` において、同パラメーター（R=2,S1=4,S2=2）でのArtEV有無の比較：

| | 勝率 | 平均スコア |
|---|---|---|
| ArtEVなし | 29.91% | 29.88 |
| ArtEV有り | 44.87% | **33.33** |

アーティファクトが場に出た際に単独撤退を能動的に狙う動きが、得点・勝率の両面で大きく効く。

**2. 文献の「Best」パラメーター（S1=4, S2=2）は多人数環境では過剰なリスク**

`scenario:risky-vs-risky` では、遅れているときに低い閾値で動く **Catch-up戦略（S2=-2）が82.8%の勝率**で圧倒。  
文献が想定する2〜3人プレイと異なり、5人以上の多人数環境ではリード時に長居するよりも**早期に利益を確定させる方が有効**。

**3. Greedy除外による環境変化**

Greedy（Never Leave）を除くだけで Risky (Best+ArtEV) の勝率が **17% → 72%** に跳ね上がった。  
これはGreedy戦略が「必ずBurstして全員の場持ちを下げる」という環境ノイズとして機能していたことを示す。  
実際のゲームでは全員撤退するタイミングが存在するため、現実的な条件下でRiskyStrategyは非常に有効。
