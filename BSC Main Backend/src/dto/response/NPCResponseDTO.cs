namespace BSC_Main_Backend.dto.response;

/// <param name="Id">Unique identifier of the NPC</param>
/// <param name="Name">Name of the NPC</param>
/// <param name="Sprite">The graphical asset used as the sprite for the NPC</param>
public record NPCResponseDTO(
    uint Id,
    string Name,
    GraphicalAssetResponseDTO Sprite  // The sprite is represented by a GraphicalAssetResponseDTO
);