namespace BSC_Main_Backend.dto.request;

/// <param name="Name">Name of the NPC</param>
/// <param name="SpriteId">The ID of the graphical asset used as the sprite for the NPC</param>
public record NPCRequestDTO(
    string Name,
    uint SpriteId  // The sprite is referenced by the ID of a GraphicalAsset
);  