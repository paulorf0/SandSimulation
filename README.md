# Sand Simulation

Simulação simples de física 2D em Go com [Ebiten](https://ebiten.org/). Blocos caem com gravidade, se empilham, se empurram e tombam quando ficam apoiados na beirada de outro bloco.

![Simulação](assets/sandsimulation.gif)

## Como rodar

```bash
go run .
```

No Linux, o Ebiten precisa das bibliotecas de desenvolvimento do X11 e OpenGL:

```bash
sudo apt install xorg-dev libgl1-mesa-dev
```

## Controles

- **Segurar o botão esquerdo do mouse**: cria blocos na posição do cursor.
