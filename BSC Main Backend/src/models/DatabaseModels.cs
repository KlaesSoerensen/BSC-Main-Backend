using System.Reflection.Metadata;

namespace BSC_Main_Backend.Models
{
    public class GraphicalAsset
    {
        public int Id { get; set; }
        public string Alias { get; set; }
        public string Type { get; set; }
        public Blob Blob { get; set; }
        public string UseCase { get; set; }
        public int Width { get; set; }
        public int Height { get; set; }
        public bool HasLOD { get; set; }
    }

    public class NPC
    {
        public int Id { get; set; }
        public int? Sprite { get; set; }
        public string Name { get; set; }
        public GraphicalAsset SpriteNavigation { get; set; }
    }

    public class Transform
    {
        public int Id { get; set; }
        public float XScale { get; set; } = 1;
        public float YScale { get; set; } = 1;
        public float XOffset { get; set; } = 0;
        public float YOffset { get; set; } = 0;
        public int ZIndex { get; set; } = 100;
    }

    public class AssetCollection
    {
        public int Id { get; set; }
        public string Name { get; set; } = "DATA.UNNAMED.COLLECTION";
        public int[] CollectionEntries { get; set; } = new int[0];
    }

    public class Achievement
    {
        public int Id { get; set; }
        public string Title { get; set; }
        public string Description { get; set; }
        public int Icon { get; set; }
        public GraphicalAsset IconNavigation { get; set; }
    }

    public class Player
    {
        public int Id { get; set; }
        public string IGN { get; set; }
        public int Sprite { get; set; }
        public int[] Achievements { get; set; } = new int[0];
        public GraphicalAsset SpriteNavigation { get; set; }
    }

    public class Session
    {
        public int Id { get; set; }
        public int Player { get; set; }
        public DateTime CreatedAt { get; set; } = DateTime.Now;
        public string Token { get; set; }
        public TimeSpan ValidDuration { get; set; } = TimeSpan.FromHours(1);
        public DateTime LastCheckIn { get; set; } = DateTime.Now;
        public Player PlayerNavigation { get; set; }
    }

    public class Colony
    {
        public int Id { get; set; }
        public string Name { get; set; } = "DATA.UNNAMED.COLONY";
        public int AccLevel { get; set; } = 0;
        public DateTime? LatestVisit { get; set; }
        public int Owner { get; set; }
        public int[] Assets { get; set; } = new int[0];
        public int[] Locations { get; set; } = new int[0];
        public int? ColonyCode { get; set; }
        public Player OwnerNavigation { get; set; }
        public ColonyCode ColonyCodeNavigation { get; set; }
    }

    public class ColonyCode
    {
        public int Id { get; set; }
        public int LobbyId { get; set; }
        public string ServerAddress { get; set; }
        public int Colony { get; set; }
        public string Value { get; set; }
        public Colony ColonyNavigation { get; set; }
    }

    public class MiniGame
    {
        public int Id { get; set; }
        public string Name { get; set; } = "DATA.UNNAMED.MINIGAME";
        public int Icon { get; set; }
        public string Description { get; set; } = "UI.DESCRIPTION_MISSING";
        public string Settings { get; set; }
        public GraphicalAsset IconNavigation { get; set; }
    }

    public class MiniGameDifficulty
    {
        public int Id { get; set; }
        public int MiniGame { get; set; }
        public int Icon { get; set; }
        public string Name { get; set; } = "?";
        public string Description { get; set; } = "UI.DESCRIPTION_MISSING";
        public string OverwritingSettings { get; set; }
        public MiniGame MiniGameNavigation { get; set; }
        public GraphicalAsset IconNavigation { get; set; }
    }

    public class Location
    {
        public int Id { get; set; }
        public string Name { get; set; } = "DATA.UNNAMED.LOCATION";
        public string Description { get; set; } = "UI.DESCRIPTION_MISSING";
        public int? MiniGame { get; set; }
        public int[] Appearances { get; set; } = new int[0];
        public MiniGame MiniGameNavigation { get; set; }
    }

    public class ColonyLocation
    {
        public int Id { get; set; }
        public int Colony { get; set; }
        public int Location { get; set; }
        public int Transform { get; set; }
        public int Level { get; set; } = 1;
        public Colony ColonyNavigation { get; set; }
        public Location LocationNavigation { get; set; }
        public Transform TransformNavigation { get; set; }
    }

    public class ColonyAsset
    {
        public int Id { get; set; }
        public int AssetCollection { get; set; }
        public int Transform { get; set; }
        public int Colony { get; set; }
        public AssetCollection AssetCollectionNavigation { get; set; }
        public Transform TransformNavigation { get; set; }
        public Colony ColonyNavigation { get; set; }
    }

    public class CollectionEntry
    {
        public int Id { get; set; }
        public int Transform { get; set; }
        public int AssetCollection { get; set; }
        public int GraphicalAsset { get; set; }
        public Transform TransformNavigation { get; set; }
        public AssetCollection AssetCollectionNavigation { get; set; }
        public GraphicalAsset GraphicalAssetNavigation { get; set; }
    }

    public class LOD
    {
        public int Id { get; set; }
        public int DetailLevel { get; set; } = 1;
        public Blob Blob { get; set; }
        public int GraphicalAsset { get; set; }
        public GraphicalAsset GraphicalAssetNavigation { get; set; }
    }
}
