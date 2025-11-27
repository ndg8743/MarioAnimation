import pygame

# 16x16 Mario sprite
MARIO = [
    "0000022222000000",
    "0000222222222000",
    "0000333113100000",
    "0003131113111000",
    "0003133111311100",
    "0003311113333000",
    "0000011111110000",
    "0000332333300000",
    "0003332332333000",
    "0033332222333300",
    "0011321221231100",
    "0011122222211100",
    "0011222222221100",
    "0000222002220000",
    "0003330000333000",
    "0033330000333300",
]

# NES-ish palette
PALETTE = {
    "0": (255, 255, 255),  # white/transparent
    "1": (240, 172, 63),   # skin/yellow
    "2": (185, 39, 22),    # red
    "3": (115, 103, 2),    # brown/green
}

PIXEL_SIZE = 20  # scale factor
WIDTH = 16 * PIXEL_SIZE
HEIGHT = 16 * PIXEL_SIZE

def main():
    pygame.init()
    screen = pygame.display.set_mode((WIDTH, HEIGHT))
    pygame.display.set_caption("Mario - Pygame Pixel Art")
    
    running = True
    while running:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                running = False
            elif event.type == pygame.KEYDOWN:
                if event.key == pygame.K_ESCAPE:
                    running = False
        
        # draw the sprite
        for y, row in enumerate(MARIO):
            for x, pixel in enumerate(row):
                color = PALETTE[pixel]
                rect = pygame.Rect(
                    x * PIXEL_SIZE,
                    y * PIXEL_SIZE,
                    PIXEL_SIZE,
                    PIXEL_SIZE
                )
                pygame.draw.rect(screen, color, rect)
        
        pygame.display.flip()
    
    pygame.quit()

if __name__ == "__main__":
    main()
