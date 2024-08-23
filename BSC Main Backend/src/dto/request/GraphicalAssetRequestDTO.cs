using System.Reflection.Metadata;

namespace BSC_Main_Backend.dto.request;

/// <param name="Alias">Alias for debugging purposes</param>
/// <param name="Type">MIME type of the asset, e.g., "image/png"</param>
/// <param name="Width">Width of the asset (resolution in pixels)</param>
/// <param name="Height">Height of the asset (resolution in pixels)</param>
/// <param name="HasLOD">Indicates whether the asset has multiple Levels of Detail (LOD)</param>
/// <param name="UseCase">Use case of the asset, e.g., "environment", "player", "icon"</param>
/// <param name="Blob">Binary data of the asset</param>
public record GraphicalAssetRequestDTO(
    string Alias,
    string Type,  // MIME type
    uint Width,
    uint Height,
    bool HasLOD,
    string UseCase,  // This could be a string that maps to AssetUseCase.ValueOf
    Blob Blob  // Optional, depending on use case
);