import pygame
import os

# ============================================
# SPRITE DATA - define your pixel art here
# ============================================

SPRITES = {
    "mario_stand": [
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
    ],
    "mario_walk1": [
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
        "0000222222220000",
        "0000022200000000",
        "0000333300000000",
        "0003333000000000",
    ],
    "mario_walk2": [
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
        "0022222222000000",
        "0000000222000000",
        "0000000333300000",
        "0000000033330000",
    ],
    "goomba": [
        "0000011111000000",
        "0001111111110000",
        "0011111111111000",
        "0111221111221100",
        "0112221111222110",
        "1112221111222111",
        "1111111111111111",
        "1111111111111111",
        "0111111111111110",
        "0011111111111100",
        "0001100000011000",
        "0011100000011100",
        "0111100000011110",
        "0111000000001110",
        "0000000000000000",
        "0000000000000000",
    ],
}

PALETTES = {
    "mario": {
        "0": (148, 148, 255),  # sky blue (transparent)
        "1": (240, 172, 63),   # skin
        "2": (185, 39, 22),    # red
        "3": (115, 103, 2),    # brown
    },
    "goomba": {
        "0": (148, 148, 255),  # sky blue (transparent)
        "1": (172, 92, 45),    # brown body
        "2": (255, 255, 255),  # white eyes
    },
}

# ============================================
# PIXEL ART ENGINE
# ============================================

class PixelSprite:
    """Convert string-based pixel art to pygame surface"""
    
    def __init__(self, data: list[str], palette: dict, scale: int = 1):
        self.data = data
        self.palette = palette
        self.width = len(data[0])
        self.height = len(data)
        self.scale = scale
        self.surface = self._render()
    
    def _render(self) -> pygame.Surface:
        """Render pixel art to a pygame surface"""
        surf = pygame.Surface(
            (self.width * self.scale, self.height * self.scale)
        )
        
        for y, row in enumerate(self.data):
            for x, pixel in enumerate(row):
                if pixel in self.palette:
                    color = self.palette[pixel]
                    rect = pygame.Rect(
                        x * self.scale,
                        y * self.scale,
                        self.scale,
                        self.scale
                    )
                    pygame.draw.rect(surf, color, rect)
        
        return surf
    
    def save(self, filename: str):
        """Save sprite as PNG"""
        pygame.image.save(self.surface, filename)
        print(f"Saved: {filename}")


class AnimatedSprite:
    """Handle sprite animation"""
    
    def __init__(self, frames: list[PixelSprite], fps: int = 8):
        self.frames = frames
        self.fps = fps
        self.current_frame = 0
        self.time_accumulator = 0
    
    def update(self, dt: float):
        """Update animation (dt in seconds)"""
        self.time_accumulator += dt
        frame_duration = 1.0 / self.fps
        
        while self.time_accumulator >= frame_duration:
            self.time_accumulator -= frame_duration
            self.current_frame = (self.current_frame + 1) % len(self.frames)
    
    @property
    def surface(self) -> pygame.Surface:
        return self.frames[self.current_frame].surface
    
    def save_spritesheet(self, filename: str):
        """Save all frames as horizontal spritesheet"""
        width = self.frames[0].surface.get_width()
        height = self.frames[0].surface.get_height()
        
        sheet = pygame.Surface((width * len(self.frames), height))
        for i, frame in enumerate(self.frames):
            sheet.blit(frame.surface, (i * width, 0))
        
        pygame.image.save(sheet, filename)
        print(f"Saved spritesheet: {filename}")


def create_sprite_from_image(image_path: str, scale: int = 1) -> pygame.Surface:
    """Load an image and scale it with nearest-neighbor (crisp pixels)"""
    img = pygame.image.load(image_path)
    if scale > 1:
        new_size = (img.get_width() * scale, img.get_height() * scale)
        img = pygame.transform.scale(img, new_size)
    return img


# ============================================
# MAIN DEMO
# ============================================

def main():
    pygame.init()
    
    SCALE = 4
    SCREEN_WIDTH = 320
    SCREEN_HEIGHT = 240
    
    screen = pygame.display.set_mode((SCREEN_WIDTH, SCREEN_HEIGHT))
    pygame.display.set_caption("Pygame Pixel Art Demo")
    clock = pygame.time.Clock()
    
    # create sprites
    mario_stand = PixelSprite(SPRITES["mario_stand"], PALETTES["mario"], SCALE)
    mario_walk1 = PixelSprite(SPRITES["mario_walk1"], PALETTES["mario"], SCALE)
    mario_walk2 = PixelSprite(SPRITES["mario_walk2"], PALETTES["mario"], SCALE)
    goomba = PixelSprite(SPRITES["goomba"], PALETTES["goomba"], SCALE)
    
    # create animation
    mario_walk = AnimatedSprite([mario_stand, mario_walk1, mario_walk2, mario_walk1], fps=8)
    
    # save outputs
    mario_stand.save("mario_stand.png")
    goomba.save("goomba.png")
    mario_walk.save_spritesheet("mario_walk_sheet.png")
    
    # positions
    mario_x = 50
    mario_y = 150
    goomba_x = 200
    goomba_y = 150
    
    # colors
    SKY_BLUE = (148, 148, 255)
    GROUND_BROWN = (139, 90, 43)
    
    running = True
    while running:
        dt = clock.tick(60) / 1000.0  # delta time in seconds
        
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                running = False
            elif event.type == pygame.KEYDOWN:
                if event.key == pygame.K_ESCAPE:
                    running = False
                elif event.key == pygame.K_s:
                    # save screenshot
                    pygame.image.save(screen, "screenshot.png")
                    print("Screenshot saved!")
        
        # update
        mario_walk.update(dt)
        
        # draw
        screen.fill(SKY_BLUE)
        
        # ground
        pygame.draw.rect(screen, GROUND_BROWN, (0, 200, SCREEN_WIDTH, 40))
        
        # sprites
        screen.blit(mario_walk.surface, (mario_x, mario_y))
        screen.blit(goomba.surface, (goomba_x, goomba_y))
        
        # instructions
        font = pygame.font.Font(None, 24)
        text = font.render("ESC=quit, S=screenshot", True, (0, 0, 0))
        screen.blit(text, (10, 10))
        
        pygame.display.flip()
    
    pygame.quit()


if __name__ == "__main__":
    main()
